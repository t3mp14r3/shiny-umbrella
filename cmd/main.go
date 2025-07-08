package main

import (
	"fmt"

	"github.com/t3mp14r3/shiny-umbrella/internal/api/handler"
	"github.com/t3mp14r3/shiny-umbrella/internal/api/usecase"
	"github.com/t3mp14r3/shiny-umbrella/internal/config"
	"github.com/t3mp14r3/shiny-umbrella/internal/logger"
	"github.com/t3mp14r3/shiny-umbrella/internal/repository"
)

func main() {
    config, err := config.New()

    if err != nil {
        return
    }

    logger, err := logger.New()
    
    if err != nil {
        return
    }

    repo, err := repository.New(config, logger)
    
    if err != nil {
        return
    }

    usecase, err := usecase.New(config, repo, logger)
    
    if err != nil {
        return
    }

    handler, err := handler.New(config, usecase, logger)
   
    if err != nil {
        return
    }

    fmt.Println(handler.Run())
}
