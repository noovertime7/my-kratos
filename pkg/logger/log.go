package logger

import (
	"backup-client/conf"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const TimeFormat = "2006-01-02 15:04:05.999"

func NewZapLogger(config *conf.Bootstrap) *zap.Logger {
	logConf := config.Log
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(TimeFormat)

	encoder := getEncoder()
	writeSyncer := getLogWriter(logConf.Filename, int(logConf.MaxSize), int(logConf.MaxBackups), int(logConf.MaxAge))
	var core zapcore.Core
	core = zapcore.NewTee(
		zapcore.NewCore(encoder, writeSyncer, cfg.Level),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), cfg.Level),
	)
	return zap.New(core)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	// 关闭caller
	encoderConfig.EncodeCaller = nil
	// 关闭message
	encoderConfig.MessageKey = ""
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
