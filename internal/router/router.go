package router

import (
	"mzt/config"
	"mzt/internal/middleware"
	"mzt/internal/service"
	"mzt/internal/validator"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authService    *service.UserService
	courseService  *service.CourseService
	paymentService *service.PaymentService
	config         *config.Config
	validator      *validator.Validator
}

func NewRouter(config *config.Config, handler *gin.Engine, authService *service.UserService, courseService *service.CourseService, paymentService *service.PaymentService, MW *middleware.Middleware) *Router {
	r := &Router{
		authService:    authService,
		paymentService: paymentService,
		courseService:  courseService,
		config:         config,
		validator:      validator.NewValidator(),
	}

	// Auth routes
	authHandler := handler.Group("/api/v1/auth")
	{

		authHandler.POST("/signin", r.SignIn)

		authHandler.POST("/signup", r.SignUp)

		authHandler.POST("/refresh", r.Refresh)
	}

	// User routes
	usersGroup := handler.Group("/api/v1/users")
	usersGroup.Use(MW.AuthMiddleware())
	{

		usersGroup.GET("/me", r.Me)
		usersGroup.GET("/me/courses", r.MyCourses)

		// Admin routes
		adminGroup := usersGroup.Group("")
		adminGroup.Use(MW.AdminVerificationMiddleware())
		{

			adminGroup.GET("/", r.GetUsers)

			adminGroup.GET("/:user_id", r.Users)

			adminGroup.PUT("/:user_id", r.Users)

			adminGroup.DELETE("/:user_id", r.Users)

			adminGroup.GET("/:user_id/role", r.Role)
		}
	}

	// Course routes
	coursesGroup := handler.Group("/api/v1/courses")

	coursesGroup.Use(MW.AuthMiddleware())
	{
		//TODO кривой слеш
		coursesGroup.GET("/", r.ListCourses)

		coursesGroup.GET("/:course_id", r.GetCourse)

		coursesGroupAdmin := coursesGroup.Group("")
		coursesGroupAdmin.Use(MW.AdminVerificationMiddleware())
		{
			coursesGroupAdmin.POST("/", r.CreateCourse)

			coursesGroupAdmin.PUT("/:course_id", r.UpdateCourse)

			coursesGroupAdmin.DELETE("/:course_id", r.DeleteCourse)
		}

		lessonsGroup := coursesGroup.Group("/:course_id/lessons")
		{
			lessonsGroup.GET("/", r.ListLessons)

			lessonsGroup.GET("/:lesson_id", r.GetLesson)

			lessonsGroupAdmin := lessonsGroup.Group("")
			lessonsGroupAdmin.Use(MW.AdminVerificationMiddleware())
			{
				lessonsGroupAdmin.POST("/", r.CreateLesson)

				lessonsGroupAdmin.PUT("/:lesson_id", r.UpdateLesson)

				lessonsGroupAdmin.DELETE("/:lesson_id", r.DeleteLesson)
			}
		}

		usersOnCourseGroup := coursesGroup.Group("/:course_id/users")
		{
			usersOnCourseGroup.POST("/", r.CreateCoursePayment)

			usersOnCourseGroup.GET("/", MW.AdminVerificationMiddleware(), r.ListUsersOnCourse)

			usersOnCourseGroup.DELETE("/:user_id", MW.AdminVerificationMiddleware(), r.RemoveUserFromCourse)
		}

		progressGroup := coursesGroup.Group("/:course_id/progress")
		{
			progressGroup.GET("/", r.GetProgress)

			progressGroup.PUT("/", r.UpdateProgress)
		}
	}

	webhookGroup := handler.Group("/api/v1/webhook/payments")
	{
		webhookGroup.POST(config.Equiring.SecretPath, r.YooWebhookHandler)
	}

	//handler.RunTLS(":443", "cert.pem", "key.pem")
	return r
}
