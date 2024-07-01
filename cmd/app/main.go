package main

import (
	"log"
	"warehousesAPI/config"
	"warehousesAPI/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("configuration error. %s", err)
	}

	app.Run(cfg)
}
