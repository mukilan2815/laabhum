package logger

import (
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps around both zap.SugaredLogger and standard log.Logger
type Logger struct {
	*zap.SugaredLogger
	*log.Logger
}

// New creates a new Logger instance with the given log level
func New(logLevel string) *Logger {
	// Create zap logger
	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig),
		zapcore.AddSync(os.Stdout),
		getZapLevel(logLevel),
	)

	zapLogger := zap.New(zapCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Create standard logger
	stdLogger := log.New(os.Stdout, "", log.LstdFlags)

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
		Logger:        stdLogger,
	}
}

// Writer returns an io.Writer for the standard Logger
func (l *Logger) Writer() io.Writer {
	return l.Logger.Writer()
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}