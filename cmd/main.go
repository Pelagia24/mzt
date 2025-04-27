package main

import (
	"mzt/config"
	"mzt/internal/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
