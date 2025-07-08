package logger

import (
	"log"

	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
    logger, err := zap.NewProduction()

    if err != nil {
        log.Printf("Failed to initialize new logger: %v\n", err)
        return nil, err
    }

    return logger, nil
}
