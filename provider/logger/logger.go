package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/robfig/cron/v3"
	"github.com/share-group/share-go/provider/config"
	"github.com/share-group/share-go/util/arrayutil"
	"github.com/share-group/share-go/util/compressutil"
	"github.com/share-group/share-go/util/fileutil"
	"github.com/share-group/share-go/util/systemutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var callerAliasMap sync.Map

type Logger struct {
	zapLogger *zap.Logger
}

var c = cron.New(cron.WithSeconds())

func init() {
	cmd, _ := os.Getwd()
	log.SetFlags(log.Flags() | log.Lshortfile | log.Lmicroseconds)
	logrotate, _ := rotatelogs.New(
		path.Join(cmd, config.GetString("logger.path"), fmt.Sprintf("%s-%s", config.GetString("application.name"), "%Y-%m-%d.log")),
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "file",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			_, file, lineNumber, _ := runtime.Caller(6)
			prefix := arrayutil.Last(strings.Split(file, "/"))
			callerAlias, _ := callerAliasMap.Load(file)
			if callerAlias != nil && len(callerAlias.(string)) > 0 {
				prefix = callerAlias.(string)
			} else {
				if !strings.HasPrefix(file, config.GetRootDir()) {
					prefix = fmt.Sprintf("share.go.%s:%d", prefix, lineNumber)
				} else {
					prefix = strings.TrimSpace(strings.ReplaceAll(file, config.GetRootDir(), ""))
					prefix = strings.TrimSpace(strings.ReplaceAll(prefix[1:], "/", "."))
					prefix = fmt.Sprintf("%s:%d", prefix, lineNumber)
				}
			}
			enc.AppendString(strings.TrimSpace(prefix))
		},
		EncodeName: zapcore.FullNameEncoder,
	}

	// 设置日志级别
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(logrotate), zapcore.AddSync(os.Stdout)}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	zap.ReplaceGlobals(zap.New(core, zap.AddCaller()))

	systemutil.Goroutine(func() { compressLogFile() })
	systemutil.Goroutine(func() { initLoggerCompress() })
}

func GetLogger(name ...string) *Logger {
	// 支持日志头重命名，默认是文件名和所在代码行数
	if len(name) > 0 && len(name[0]) > 0 {
		stackArray := strings.Split(string(debug.Stack()), "\n")
		caller := strings.TrimSpace(stackArray[6])
		caller = strings.TrimSpace(caller[0:strings.LastIndex(caller, ":")])
		callerAliasMap.Store(caller, strings.TrimSpace(name[0]))
	}
	return &Logger{zapLogger: zap.L()}
}

func (o *Logger) Debug(msg string, args ...any) {
	o.zapLogger.Debug(fmt.Sprintf(msg, args...))
}

func (o *Logger) Info(msg string, args ...any) {
	o.zapLogger.Info(fmt.Sprintf(msg, args...))
}

func (o *Logger) Warn(msg string, args ...any) {
	o.zapLogger.Warn(fmt.Sprintf(msg, args...))
}

func (o *Logger) Error(msg string, args ...any) {
	o.zapLogger.Error(fmt.Sprintf(msg, args...))
}

func (o *Logger) Panic(msg string, args ...any) {
	o.zapLogger.Panic(fmt.Sprintf(msg, args...))
}

func (o *Logger) Fatal(msg string, args ...any) {
	o.zapLogger.Fatal(fmt.Sprintf(msg, args...))
	os.Exit(0)
}

func initLoggerCompress() {
	c.AddFunc("1 0 0 * * *", func() {
		compressLogFile()
	})
	c.Start()
}

func compressLogFile() {
	suffix := ".tar.bz2"
	location, _ := time.LoadLocation("Local")
	yesterday, _ := time.ParseInLocation(time.DateOnly, time.Now().AddDate(0, 0, -1).Format(time.DateOnly), time.Now().In(location).Location())
	regex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	loggerPath := strings.TrimSpace(path.Join(config.GetRootDir(), config.GetString("logger.path")))
	for _, logFile := range fileutil.ListDir(loggerPath) {
		if strings.HasSuffix(logFile, suffix) {
			continue
		}

		date := strings.TrimSpace(regex.FindString(logFile))
		thisDay, _ := time.ParseInLocation(time.DateOnly, date, time.Now().In(location).Location())
		if thisDay.UnixMilli() > yesterday.UnixMilli() {
			continue
		}

		compressFile := path.Join(loggerPath, fmt.Sprintf("%s-%s%s", config.GetString("application.name"), date, suffix))
		if fileutil.Exists(compressFile) {
			continue
		}

		compressutil.Bzip2Compress(compressFile, logFile)
		os.RemoveAll(logFile)
	}
}
