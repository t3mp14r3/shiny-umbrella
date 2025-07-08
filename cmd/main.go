package main

import (
	"fmt"

	"github.com/t3mp14r3/shiny-umbrella/internal/config"
	"github.com/t3mp14r3/shiny-umbrella/internal/logger"
)

func main() {
    _, err := config.New()

    if err != nil {
        return
    }

    _, err = logger.New()
    
    if err != nil {
        return
    }

    fmt.Println("all good")
}
