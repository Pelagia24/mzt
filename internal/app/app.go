package app

import (
	"mzt/config"
	"mzt/internal/middleware"
	"mzt/internal/migration"
	"mzt/internal/repository"
	"mzt/internal/router"
	"mzt/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// создаем репозитории для работы с данными
	userRepo := repository.NewUserRepo(cfg)
	courseRepo := repository.NewCourseRepo(cfg)
	eventRepo := repository.NewEventRepo(cfg)
	paymentRepo := repository.NewPaymentRepo(cfg)

	// запускаем миграции базы данных
	migration.RunMigrations(cfg)

	// создаем сервисы для бизнес логики
	authService := service.NewUserService(cfg, userRepo)
	courseService := service.NewCourseService(cfg, courseRepo)
	paymentService := service.NewPaymentService(cfg, courseRepo, paymentRepo)
	eventService := service.NewEventService(cfg, eventRepo, courseRepo)

	// создаем middleware для обработки запросов
	middleware := middleware.NewMiddleware(cfg, userRepo, courseRepo)

	// создаем роутер
	handler := gin.Default()

	// настраиваем cors для работы с фронтендом
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080", "http://127.0.0.1:5173", "http://127.0.0.1:8080", "https://c221-62-60-236-43.ngrok-free.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// настраиваем все маршруты
	router.NewRouter(cfg, handler, authService, courseService, paymentService, eventService, middleware)
	// запускаем сервер на порту 8080
	handler.Run(":8080")
	//TODO server
}
