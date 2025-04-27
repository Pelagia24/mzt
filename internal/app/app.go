package app

import (
	"mzt/config"
	"mzt/internal/middleware"
	"mzt/internal/repository"
	"mzt/internal/router"
	"mzt/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	userRepo := repository.NewUserRepo(cfg)
	courseRepo := repository.NewCourseRepo(cfg)

	Migrate(userRepo)

	authService := service.NewUserService(cfg, userRepo)
	courseService := service.NewCourseService(cfg, courseRepo)

	middleware := middleware.NewMiddleware(cfg, userRepo)

	handler := gin.Default()

	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	router.NewRouter(cfg, handler, authService, courseService, middleware)
	handler.Run(":8080")
	//TODO server
}
