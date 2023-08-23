// logger package handles logging utility in the application
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(logDir string) *zap.Logger {
	env := os.Getenv("ENVIRONMENT_NAME")
	opts := []zap.Option{
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if env == "prod" || env == "production" {
		return zap.New(fileCore(logDir), opts...)
	}
	return zap.New(
		zapcore.NewTee(
			fileCore(logDir),
			consoleCore(),
		),
		opts...,
	)
}

// fileCore returns core for log file
func fileCore(logDir string) zapcore.Core {
	encodeConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "Name",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stack",
		SkipLineEnding: false,
		LineEnding:     "",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encodeConfig)

	now := time.Now()

	formattedDate := fmt.Sprintf("%02d-%s-%d", now.Day(), now.Month().String()[:3], now.Year())

	logFile, err := createEmptyLogFile(logDir, formattedDate)
	if err != nil {
		panic(fmt.Sprintf("could not create empty log file : %s", err))
	}

	rw := rotateWriter{
		filename: logFile.Name(),
		fp:       logFile,
	}

	env := os.Getenv("ENVIRONMENT_NAME")
	var level zapcore.LevelEnabler

	if (env != "prod") || (env != "production") {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(&rw), level)
	return zapcore.RegisterHooks(core, func(e zapcore.Entry) error {
		t := e.Time
		currentDate := fmt.Sprintf("%02d-%s-%d", t.Day(), t.Month().String()[:3], t.Year())

		// if current date is same as logFile name then do nothing
		if strings.Contains(logFile.Name(), currentDate) {
			return nil
		}
		rw.filename = logDir + "/" + currentDate + ".log"
		return rw.Rotate()
	})
}

// consoleCore returns the core for console log
func consoleCore() zapcore.Core {
	encodeConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "Name",
		CallerKey:      "Caller",
		FunctionKey:    "Function",
		StacktraceKey:  "trace",
		SkipLineEnding: false,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encodeConfig)
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
}

func createEmptyLogFile(directoryPath, filename string) (*os.File, error) {
	if err := createDir(directoryPath); err != nil {
		return nil, err
	}
	filePath := directoryPath + "/" + filename + ".log"
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// createDir crates folder if it does not exist
func createDir(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirname, 0o755)
		if err != nil {
			return errDir
		}
	}
	return nil
}
