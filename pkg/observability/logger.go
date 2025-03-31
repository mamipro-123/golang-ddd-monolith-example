package observability

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// InitLogger initializes the global logger with the specified environment
func InitLogger(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout", "stderr"}
		config.Encoding = "json"
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}
		config.Encoding = "console"
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	var err error
	globalLogger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// Fallback to a basic logger if not initialized
		globalLogger, _ = zap.NewProduction()
	}
	return globalLogger
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
} 