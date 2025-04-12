package app

import (
	"mzt/config"
	"mzt/internal/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {

	repo := auth.NewRefreshTokensRepo(cfg)

	Migrate(repo)

	service := auth.NewService(cfg, repo)

	middleware := auth.NewMiddleware(cfg, repo)

	handler := gin.Default()

	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	auth.NewRouter(cfg, handler, service, middleware)
	handler.Run(":8080")
	//TODO server
	//TODO nginx
}
