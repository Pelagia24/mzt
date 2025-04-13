package auth

import (
	"mzt/config"
	_ "mzt/docs"
	"mzt/internal/auth/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	usersGroup := handler.Group("/users")

	usersGroup.Use(MW.AuthMiddleware())
	{
		usersGroup.GET("/me", r.Me)

		adminGroup := usersGroup.Group("")
		adminGroup.GET("/", r.GetUsers)
		adminGroup.Use(MW.AdminVerificationMiddleware()).Any("/:user_id", r.Users)
		adminGroup.GET("/:user_id/role", r.Role)
	}

	//handler.RunTLS(":443", "cert.pem", "key.pem")
	return r
}

func (r *Router) GetUsers(c *gin.Context) {
	users, err := r.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User infos",
		"users":   users,
	})
}

func (r *Router) Role(c *gin.Context) {
	userId := c.Param("user_id")
	id, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	role, err := r.service.Role(id)
	if err != nil || role == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't get user role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Role of user",
		"role":    role,
	})
}

func (r *Router) Users(c *gin.Context) {
	userId := c.Param("user_id")
	id, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	switch c.Request.Method {
	case http.MethodGet:
		user, err := r.service.GetUser(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Can't get user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User detailedinfo",
			"user":    user,
		})

	case http.MethodDelete:
		self, ok := c.Get("self")
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unknown sender"})
			return
		}
		casted, ok := self.(string)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unknown sender"})
		}

		selfId, err := uuid.Parse(casted)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't parse user id"})
			return
		}

		if id == selfId {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't delete self"})
			return
		}

		err = r.service.DeleteUser(id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deleted user with success"})
	case http.MethodPut:
		var payload dto.UpdateUserDto
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := r.service.UpdateUser(id, &payload)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Updated user with success"})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid http method"})
		return
	}

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
	var payload dto.LoginDto
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
	var payload dto.RegistrationDto
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
