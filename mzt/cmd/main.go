package main

import (
	"mzt/config"
	_ "mzt/docs"
	"mzt/internal/auth/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
