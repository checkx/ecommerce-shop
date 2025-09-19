package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"ecommerce-shop/internal/config"
)

func New(cfg config.Config) *zap.Logger {
	var z *zap.Logger
	if cfg.Env == "production" {
		cfgZap := zap.NewProductionConfig()
		cfgZap.OutputPaths = []string{"stdout"}
		cfgZap.ErrorOutputPaths = []string{"stderr"}
		z, _ = cfgZap.Build()
	} else {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		z = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
	}
	return z
}
