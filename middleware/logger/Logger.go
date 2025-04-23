package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logrus.SetLevel(logrus.InfoLevel)
	maxSize, errS := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
	maxBackup, errB := strconv.Atoi(os.Getenv("LOG_MAX_BACKUP"))
	maxAge, errA := strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
	isCompress, errC := strconv.ParseBool(os.Getenv("LOG_MAX_COMPRESS"))

	if errS != nil || errB != nil || errA != nil || errC != nil {
		panic("Error during conversion")
	}

	logFile := &lumberjack.Logger{
		Filename:   os.Getenv("LOG_FILE"),
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   isCompress,
	}

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(io.MultiWriter(logFile, os.Stdout))
	return logger
}
