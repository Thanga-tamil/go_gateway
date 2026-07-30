package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func Init(logPath string, logLevel int) {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zapcore.Level

	switch logLevel {
	case 1:
		level = zap.InfoLevel
	case 2:
		level = zap.ErrorLevel
	case 3:
		level = zap.DebugLevel
	default:
		level = zap.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(file), zapcore.AddSync(os.Stdout)),
		level,
	)

	logger := zap.New(core)
	Log = logger.Sugar()
}


func Info(args ...interface{}) {
	Log.Info(append([]interface{}{"ℹ️"}, args...)...)
}

func Infof(template string, args ...interface{}) {
	Log.Infof("ℹ️ "+template, args...)
}

func Error(args ...interface{}) {
	Log.Error(append([]interface{}{"❌"}, args...)...)
}

func Errorf(template string, args ...interface{}) {
	Log.Errorf("❌ "+template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Log.Fatalf("💥 "+template, args...)
}

func Debug(args ...interface{}) {
	Log.Debug(append([]interface{}{"ℹ️"}, args...)...)
}

