package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/share-group/share-go/provider/config"
	"github.com/share-group/share-go/util/arrayutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const frameworkName = "share-go"

var callerAliasMap sync.Map

func init() {
	cmd, _ := os.Getwd()
	log.SetFlags(log.Flags() | log.Lshortfile)
	logrotate, _ := rotatelogs.New(
		path.Join(cmd, config.GetString("logger.path")+"/"+config.GetString("application.name")+"-%Y-%m-%d.log"),
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
			prefix := strings.ReplaceAll(caller.FullPath(), strings.ReplaceAll(cmd, "\\", "/"), "")
			callerAlias, _ := callerAliasMap.Load(caller.FullPath()[0:strings.LastIndex(caller.FullPath(), ":")])
			if callerAlias != nil && len(callerAlias.(string)) > 0 {
				prefix = callerAlias.(string)
			} else {
				if strings.Contains(prefix, frameworkName) {
					atvIndex := strings.Index(strings.ReplaceAll(prefix, frameworkName, " "), "@v")
					if atvIndex > -1 {
						prefix = prefix[atvIndex:]
						prefix = prefix[strings.Index(prefix, "/"):]
						prefix = strings.ReplaceAll(prefix[1:], "/", ".")
					} else {
						prefix = fmt.Sprintf("%s", strings.ReplaceAll(prefix[strings.Index(prefix, frameworkName)+len(frameworkName)+1:], "/", "."))
					}
					prefix = prefix[:strings.LastIndex(prefix, ".")]
					prefix = arrayutil.Last(strings.Split(prefix, "."))
					prefix = fmt.Sprintf("share.go.%s", prefix)
				} else {
					prefix = strings.ReplaceAll(prefix[1:], "/", ".")
					lastDotIndex := strings.LastIndex(prefix, ".")
					prefix = fmt.Sprintf("%s.%s%s", getFirstLetter(prefix[:lastDotIndex]), arrayutil.Last(strings.Split(caller.Function, ".")), strings.ReplaceAll(prefix[lastDotIndex:], ".go", ""))
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
}

func GetLogger(name ...string) *zap.Logger {
	// 支持日志头重命名，默认是文件名和所在代码行数
	if len(name) > 0 && len(name[0]) > 0 {
		stackArray := strings.Split(string(debug.Stack()), "\n")
		caller := strings.TrimSpace(stackArray[6])
		caller = caller[0:strings.LastIndex(caller, ":")]
		callerAliasMap.Store(caller, name[0])
	}
	return zap.L()
}

func getFirstLetter(str string) string {
	result := make([]string, 0)
	for _, s := range strings.Split(str, ".") {
		result = append(result, s[:1])
	}
	return strings.Join(result, ".")
}
