package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 創建 logger
//   - fields: 固定欄位列表
func NewLogger(config *Config, fields zap.Option) *zap.Logger {

	encConfig := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "Level",
		MessageKey:     "Msg",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.999Z07:00"),
		EncodeDuration: zapcore.MillisDurationEncoder,
	}
	enc := zapcore.NewJSONEncoder(encConfig)

	// logrotate: https://github.com/uber-go/zap/blob/master/FAQ.md
	writeSyncers := []zapcore.WriteSyncer{}
	if config.Enable {
		if len(config.FilenPath) > 0 {
			wsFile := zapcore.AddSync(&lumberjack.Logger{
				Filename:   config.FilenPath + config.Filename,
				MaxSize:    config.MaxSize,
				MaxAge:     config.MaxAge,
				MaxBackups: config.MaxBackups,
			})
			writeSyncers = append(writeSyncers, wsFile)
		}
		wsStdout := zapcore.AddSync(os.Stdout)
		writeSyncers = append(writeSyncers, wsStdout)
	}

	core := zapcore.NewCore(enc, zapcore.NewMultiWriteSyncer(writeSyncers...), config.getZapLevel())

	options := []zap.Option{}
	if fields != nil {
		options = append(options, fields)
	}

	zapLogger := zap.New(core, options...)
	return zapLogger
}
