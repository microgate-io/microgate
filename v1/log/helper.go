package log

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap/zapcore"
)

func Init() {
	// logger, _ := zap.NewProduction()
	cfg := zapdriver.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Encoding = "console"
	logger, _ := cfg.Build()
	defer logger.Sync()
	InitLogger(logger)
}
