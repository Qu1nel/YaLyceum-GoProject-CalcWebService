package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"CalcService/config"
	"CalcService/internal/server"
	"CalcService/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Init logging and context
	mainLogger := logger.New()
	ctx := context.WithValue(context.Background(), logger.LoggerKey, mainLogger)

	cfg, err := config.New() // Config instance
	if err != nil {
		mainLogger.Error(ctx, config.LogErrorLoadConfig, zap.String(config.LogErrorMsg, err.Error())) // Ошибка загрузки конфига
		return
	}

	// Main app
	server, err := server.New(ctx, cfg.RestServerPort, cfg.PatternURL)
	if err != nil {
		mainLogger.Error(ctx, config.LogErrorInitgRPCServer, zap.String(config.LogErrorMsg, err.Error())) // Ошибка запуска gRPC сервера
		return
	}
	// Context for shutdown
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Get port from config
	port := fmt.Sprintf(":%d", cfg.RestServerPort)

	// Start server listening..
	go func() {
		if err := server.Start(port); err != nil && err != http.ErrServerClosed {
			mainLogger.Error(ctx, config.LogShutdownServer, zap.String(config.LogErrorCtx, err.Error())) // Ошибка запуска сервера
			return
		}
	}()
	mainLogger.Info(ctx, config.LogInitServer, zap.Int("port", cfg.RestServerPort)) // Запуск сервера

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown of server
	if err := server.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
		return
	}

	mainLogger.Info(ctx, config.LogServerStopped)
}
