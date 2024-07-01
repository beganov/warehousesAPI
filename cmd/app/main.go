package main

import (
	"github.com/robertgarayshin/warehousesAPI/config"
	"github.com/robertgarayshin/warehousesAPI/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("configuration error. %s", err)
	}

	app.Run(cfg)
}
