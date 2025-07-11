package logger

import (
	"log"

	"go.uber.org/zap"
)

func New() *zap.Logger {
    logger, err := zap.NewProduction()

    if err != nil {
        log.Fatalf("Failed to initialize new logger: %v\n", err)
    }

    return logger
}
