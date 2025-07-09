package main

import (
	"fmt"

	"github.com/t3mp14r3/shiny-umbrella/internal/api/handler"
	"github.com/t3mp14r3/shiny-umbrella/internal/api/usecase"
	"github.com/t3mp14r3/shiny-umbrella/internal/config"
	"github.com/t3mp14r3/shiny-umbrella/internal/cron"
	"github.com/t3mp14r3/shiny-umbrella/internal/logger"
	"github.com/t3mp14r3/shiny-umbrella/internal/repository"
	"github.com/t3mp14r3/shiny-umbrella/internal/repository/notifier"
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

    cron, err := cron.New(repo, logger)
    
    if err != nil {
        return
    }

    err = cron.Load()
    
    if err != nil {
        return
    }

    notifier, err := notifier.New(logger, cron, config.RepositoryConnString())
    
    if err != nil {
        return
    }
    
    go notifier.Listen()

    usecase, err := usecase.New(repo, logger)
    
    if err != nil {
        return
    }

    handler, err := handler.New(config, usecase, logger)
   
    if err != nil {
        return
    }

    fmt.Println(handler.Run())
}
