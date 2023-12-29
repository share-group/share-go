package bootstrap

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/share-group/share-go/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const frameworkName = "share-go"

type _logger struct{}

var Logger = newLogger()

func newLogger() *_logger {
	return &_logger{}
}

func init() {
	cmd, _ := os.Getwd()
	log.SetFlags(log.Flags() | log.Lshortfile)
	logrotate, _ := rotatelogs.New(
		path.Join(cmd, Config.GetStringValue("logger.path")+"/"+Config.GetStringValue("application.name")+"-%Y-%m-%d.log"),
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
			if strings.Contains(prefix, frameworkName) {
				prefix = fmt.Sprintf("%s", strings.ReplaceAll(prefix[strings.Index(prefix, frameworkName)+len(frameworkName)+1:], "/", "."))
			} else {
				prefix = strings.ReplaceAll(prefix[1:], "/", ".")
			}
			prefix = strings.ReplaceAll(prefix, ".go:", fmt.Sprintf(".%s:", util.ArrayUtil.Last(strings.Split(caller.Function, "."))))
			enc.AppendString(prefix)
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

func (l *_logger) GetLogger() *zap.Logger {
	return zap.L()
}
