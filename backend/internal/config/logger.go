package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg Config) *zap.Logger {
	zcfg := zap.NewProductionConfig()
	if cfg.Server.Mode == "debug" {
		zcfg = zap.NewDevelopmentConfig()
	}
	zcfg.EncoderConfig.TimeKey = "ts"
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zcfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
