package logger

import (
	"context"
	"log"

	"go.uber.org/zap"
)

type LoggerKeyType string

const (
	LoggerKey LoggerKeyType = "logger"
	RequestID string        = "requestID"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}
type logger struct {
	logger *zap.Logger
}

// Info
func (l logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.logger.Info(msg, fields...)
}

// Error
func (l logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.logger.Error(msg, fields...)
}

func New() Logger {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer zapLogger.Sync() //nolint:errcheck
	return &logger{
		logger: zapLogger,
	}
}

func GetLoggerFromCtx(ctx context.Context) Logger {
	return ctx.Value(LoggerKey).(Logger)
}
