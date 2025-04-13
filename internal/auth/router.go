package auth

import (
	"fmt"
	"mzt/config"
	_ "mzt/docs"
	"mzt/internal/auth/dto"
	"mzt/internal/auth/utils"
	"net/http"
	"regexp"

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

	usersGroup := handler.Group("api/v1/users")

	usersGroup.Use(MW.AuthMiddleware())
	{
		usersGroup.GET("/me", r.Me)

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

	//handler.RunTLS(":443", "cert.pem", "key.pem")
	return r
}

// @Summary Get all users info(only admin)
// @Description Gets all users
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/ [get]
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

// @Summary Get user role by id
// @Description Gets all users
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/:user_id/role [get]
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

// @Summary Update user by dto (only admin)
// @Description Updates user
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UpdateUserDto true "User to update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/:user_id [put]
func (r *Router) _UsersPut(c *gin.Context) {}

// @Summary Get user by dto (only admin)
// @Description Gets user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/:user_id [get]
func (r *Router) _UsersGet(c *gin.Context) {}

// @Summary Delete user by dto (only admin)
// @Description Deletes user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/:user_id [delete]
func (r *Router) _UsersDelete(c *gin.Context) {}

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
// @Param user body dto.LoginDto true "User login info"
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

	if !isValidEmail(payload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	if !isValidPassword(payload.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is weak or contains incorrect symbols"})
		return
	}

	access, refresh, err := r.service.SignIn(&payload)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := r.service.GetUserId(payload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	role, err := r.service.Role(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", true, true)
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "User signed in successfully",
		"access_token": access,
		"id":           id,
		"role":         role,
	})
}

// @Summary Sign up user
// @Description Sign up user
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.RegistrationDto true "User to create"
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

	if !isValidEmail(payload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
	}

	if !isValidPassword(payload.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
	}
	if !isValidName(payload.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
	}
	if !isValidPhoneNumber(payload.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
	}
	if !isValidTelegram(payload.Telegram) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram"})
	}

	access, refresh, err := r.service.SignUp(&payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := r.service.GetUserId(payload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	role, err := r.service.Role(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", true, true)
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "User created successfully",
		"access_token": access,
		"id":           id,
		"role":         role,
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
	fmt.Println("1")
	access, refresh, err := r.service.RefreshTokens(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//after validation
	parsed, err := utils.ValidateToken(token, r.config.Jwt.RefreshKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("2")

	sub, err := parsed.Claims.GetSubject()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("3")

	id, err := r.service.GetUserId(sub)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("4")

	role, err := r.service.Role(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", true, true)
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Tokens refreshed successfully",
		"access_token": access,
		"id":           id,
		"role":         role,
	})

}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func isValidTelegram(telegram string) bool {
	const telegramRegex = `^@[A-Za-z0-9_]{5,32}$`
	re := regexp.MustCompile(telegramRegex)
	return re.MatchString(telegram)
}

func isValidPhoneNumber(phone string) bool {
	const phoneRegex = `^\+?\d{10,15}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

func isValidPassword(passwd string) bool {
	if len(passwd) < 8 {
		return false
	}

	allowedChars := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()\[\]\-_=+{}|;:,.<>?/]+$`).MatchString(passwd)
	return allowedChars
}

func isValidName(name string) bool {
	validName := regexp.MustCompile(`^[A-Za-z]{2,}$`)
	return validName.MatchString(name)
}
