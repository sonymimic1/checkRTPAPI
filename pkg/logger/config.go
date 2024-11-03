package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 設定
type Config struct {
	Enable     bool   `json:"enable"` // 是否啟用
	FilenPath  string // log 檔路徑; 若為空字串, 則不輸出到檔案
	Filename   string // log 檔名
	MaxSize    int    // log rotate 的檔案大小 (MB)
	MaxAge     int    // log 檔名日期的保留天數
	MaxBackups int    // log 檔的保留數量
	Level      string `json:"level"` // log level
}

// level 列表 (levelName -> zapLevel)
var levelMap = map[string]zapcore.Level{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"dpanic": zap.DPanicLevel,
	"panic":  zap.PanicLevel,
	"fatal":  zap.FatalLevel,
}

//#region private method

// 取得 zap level
func (c *Config) getZapLevel() zapcore.Level {
	level := zap.InfoLevel
	if zapLevel, ok := levelMap[c.Level]; ok {
		level = zapLevel
	}
	return level
}

//#endregion
