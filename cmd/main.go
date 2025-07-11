package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/t3mp14r3/shiny-umbrella/internal/api/handler"
	"github.com/t3mp14r3/shiny-umbrella/internal/api/usecase"
	"github.com/t3mp14r3/shiny-umbrella/internal/config"
	"github.com/t3mp14r3/shiny-umbrella/internal/cron"
	"github.com/t3mp14r3/shiny-umbrella/internal/logger"
	"github.com/t3mp14r3/shiny-umbrella/internal/repository"
	"github.com/t3mp14r3/shiny-umbrella/internal/repository/notifier"
	"go.uber.org/zap"
)

func main() {
    config := config.New()

    logger := logger.New()

    repo := repository.New(config, logger)
    defer repo.Close()

    cron := cron.New(repo, logger)
    defer cron.Close()
    
    notifier := notifier.New(logger, cron, config.RepositoryConnString())
    defer notifier.Close()
    
    usecase := usecase.New(repo, logger)
    
    handler := handler.New(config, usecase, logger)
   
    wg := sync.WaitGroup{}
    ctx, cancel := context.WithCancel(context.Background())

    cron.Start()

    err := cron.Load()
    
    if err != nil {
        logger.Error("Failed to execute cron load", zap.Error(err))
        return
    }

    wg.Add(1)
    go func(ctx context.Context) {
        defer wg.Done()
        if err := notifier.Listen(ctx); err != nil {
            logger.Error("Notifier error", zap.Error(err))
            cancel()
        }
    }(ctx)

    wg.Add(1)
    go func(ctx context.Context) {
        defer wg.Done()
        if err := handler.Run(ctx); err != nil {
            logger.Error("Handler error", zap.Error(err))
            cancel()
        }
    }(ctx)

    logger.Info("Application started", zap.String("addr", config.AppAddr))

    exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

    select {
        case <-ctx.Done():
            logger.Error("Stopping via context")
        case <-exit:
            logger.Info("Stopping")
    }

    cancel()
    wg.Wait()

    logger.Info("Application stopped")
}
