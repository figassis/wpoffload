package util

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//cfg.DisableStacktrace = true
	return cfg.Build()
}

// Simple log function for logging to a file
func Log(logmessage interface{}) (newerr error) {

	level := 3
	message := fmt.Sprint(logmessage)
	newerr = nil
	if logmessage == nil {
		return
	}

	switch msgType := logmessage.(type) {
	case string:
		message = msgType
		level = 1
	case error:
		if msgType == nil {
			return
		}
		level = 3
		newerr = msgType
	}

	switch level {
	case 3, 4:
		logger.Error(message, zap.String("version", Version))
	case 2:
		logger.Warn(message, zap.String("version", Version))
	case 1:
		logger.Info(message, zap.String("version", Version))
	default:
		logger.Debug(message, zap.String("version", Version))
	}

	return
}
