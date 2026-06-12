package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if err := config.Level.UnmarshalText([]byte(strings.ToLower(strings.TrimSpace(level)))); err != nil {
		config.Level.SetLevel(zap.InfoLevel)
	}
	return config.Build()
}
