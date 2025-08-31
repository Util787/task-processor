package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Util787/task-processor/internal/adapters/rest"
	"github.com/Util787/task-processor/internal/config"
	"github.com/Util787/task-processor/internal/infra/storage"
	taskprocessqueue "github.com/Util787/task-processor/internal/infra/task-process-queue"
	"github.com/Util787/task-processor/internal/usecase"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)

	// storage
	inMemStorage := storage.NewInMemoryTaskStateStorage()

	// task process queue
	log.Info("Starting task process queue", slog.Int("num_of_workers", cfg.TaskProcessQueueConfig.Workers), slog.Int("queue_size", cfg.TaskProcessQueueConfig.QueueSize))
	taskProcessQueue := taskprocessqueue.NewTaskProcessQueue(context.Background(), log, cfg.TaskProcessQueueConfig, inMemStorage)

	// usecase
	taskUsecase := usecase.NewTaskUsecase(inMemStorage, taskProcessQueue)

	// server
	server := rest.NewHTTPServer(log, cfg.HTTPServerConfig, taskUsecase)

	// start
	go func() {
		log.Info("Starting server", slog.String("host", cfg.HTTPServerConfig.Host), slog.Int("port", cfg.HTTPServerConfig.Port))
		err := server.Run()
		if err != nil {
			log.Error("Server interrupted", slog.String("error", err.Error()))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Info("Shutting down gracefully...")

	log.Info("Shutting down server")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("Server shutdown error", slog.String("error", err.Error()))
	}

	log.Info("Shutting down task process queue")
	taskProcessQueue.Shutdown()

	log.Info("Shutdown complete")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
