// Package logger provides a zap-based logger implementation.
// It supports both JSON and console output formats with configurable log levels.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultLogger creates a default zap logger configuration.
// It supports both JSON and console output formats and configurable log levels.
//
// Parameters:
//   - level: The log level (-1=Debug, 0=Info, 1=Warn, 3=DPanic, 4=Panic, 5=Fatal, default=Error)
//   - format: The output format ("console" for human-readable, "json" for structured logging)
//
// Returns:
//   - *zap.Logger: A configured zap logger instance
func DefaultLogger(level int, format string) *zap.Logger {
	var coreArr []zapcore.Core

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // Specify time format
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder  // Capital level encoder (not color for JSON)
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // Display short file path

	var encoder zapcore.Encoder
	// NewJSONEncoder() outputs JSON format, NewConsoleEncoder() outputs plain text format
	if format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Use color for console
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Define log level enablers
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		// Error level and above
		return lev >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		// Info and debug levels (debug is the lowest)
		return lev < zapcore.ErrorLevel && lev >= zapcore.DebugLevel
	})

	// Create cores for different priority levels
	infoCore := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		lowPriority,
	)
	errorCore := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		highPriority,
	)

	coreArr = append(coreArr, infoCore)
	coreArr = append(coreArr, errorCore)

	// Create logger with cores and set caller info
	logger := zap.New(
		zapcore.NewTee(coreArr...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return setLogLevel(logger, level)
}

// setLogLevel configures the minimum log level for the logger.
//
// Parameters:
//   - logger: The zap logger to configure
//   - level: The minimum log level (-1=Debug, 0=Info, 1=Warn, 3=DPanic, 4=Panic, 5=Fatal, default=Error)
//
// Returns:
//   - *zap.Logger: The logger with the configured level
func setLogLevel(logger *zap.Logger, level int) *zap.Logger {
	var zapLevel zapcore.Level

	switch level {
	case -1:
		zapLevel = zapcore.DebugLevel
	case 0:
		zapLevel = zapcore.InfoLevel
	case 1:
		zapLevel = zapcore.WarnLevel
	case 3:
		zapLevel = zapcore.DPanicLevel
	case 4:
		zapLevel = zapcore.PanicLevel
	case 5:
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.ErrorLevel
	}

	return logger.WithOptions(zap.IncreaseLevel(zapLevel))
}
