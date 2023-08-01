package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var Logger *log.Logger

func init() {
	Logger = log.New()
	cmd, _ := os.Getwd()
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&log.TextFormatter{})

	logf, err := rotatelogs.New(
		path.Join(cmd, "log/app.%Y%m%d.log"),
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	Logger.SetOutput(logf)
	if err != nil {
		log.Printf("failed to create rotatelogs: %s", err)
		return
	}
}
