package main

import (
	"log"
	"os"

	"github.com/teamlix/user-service/internal/app"
)

func main() {
	configPath := "config/config.yaml"
	err := app.Run(configPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
