package auth

import (
	"mzt/config"
	_ "mzt/docs"
	"mzt/internal/auth/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	service *Service
	config  *config.Config
}

// @title API
// @version 1.0
// @description This is an API for authentication and authorization.
// @contact.name API Support
// @contact.email ABOBA
// @license.name MIT
// @license.url http://opensource.org/licenses/MIT
// @host localhost:8080
// @basePath /api/v1
func NewRouter(config *config.Config, handler *gin.Engine, service *Service, MW *Middleware) *Router {
	r := &Router{
		service: service,
		config:  config,
	}

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authHandler := handler.Group("/api/v1/auth")
	{
		authHandler.POST("/signin", r.SignIn)
		authHandler.POST("/signup", r.SignUp)
		authHandler.POST("/refresh", r.Refresh)
	}

	secureHandler := handler.Group("/api/v1/users")
	secureHandler.Use(MW.AuthMiddleware())
	{
		secureHandler.GET("/me", r.Me)
		//secureHandler.POST("/logout", r.Logout)
	}
	//handler.RunTLS(":443", "cert.pem", "key.pem")
	return r
}

// @Summary Get user profile
// @Description Get user profile
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/me [get]
func (r *Router) Me(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile fetched successfully",
		"user":    user,
	})
}

// @Summary Sign in user
// @Description Sign in user
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserLoginInfo true "User login info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/signin [post]
func (r *Router) SignIn(c *gin.Context) {
	var payload dto.UserLoginInfo
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := r.service.SignIn(&payload)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", true, true)
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "User signed in successfully",
		"access_token": access,
	})
}

// @Summary Sign up user
// @Description Sign up user
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.User true "User to create"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/signup [post]
func (r *Router) SignUp(c *gin.Context) {
	var payload dto.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := r.service.SignUp(&payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", true, true)
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "User created successfully",
		"access_token": access,
	})
}

// @Summary Refresh tokens
// @Description Refresh tokens
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (r *Router) Refresh(c *gin.Context) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := r.service.RefreshTokens(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", true, true)
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Tokens refreshed successfully",
		"access_token": access,
	})

}
