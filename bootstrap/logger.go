package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strings"
	"time"
)

func init() {
	cmd, _ := os.Getwd()
	logrotate, _ := rotatelogs.New(
		path.Join(cmd, fmt.Sprintf("%v", config.Get("logger.path"))+"/"+fmt.Sprintf("%v", config.Get("application.name"))+"-%Y-%m-%d.log"),
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
			prefix = strings.ReplaceAll(prefix[1:], "/", "|")
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
