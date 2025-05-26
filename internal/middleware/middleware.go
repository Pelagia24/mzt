package middleware

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/repository"
	"mzt/internal/service"
	"mzt/internal/validator"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Middleware struct {
	config     *config.Config
	repo       *repository.UserRepo
	courseRepo *repository.CourseRepo
	validator  *validator.Validator
}

func NewMiddleware(config *config.Config, repo *repository.UserRepo, courseRepo *repository.CourseRepo) *Middleware {
	return &Middleware{
		config:     config,
		repo:       repo,
		courseRepo: courseRepo,
		validator:  validator.NewValidator(),
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(AuthHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect authorization"})
			return
		}

		tokenString := fields[1]
		token, err := m.validator.ValidateToken(tokenString, m.config.Jwt.AccessKey)

		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil || sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Corrupted sub"})
			return
		}

		user, err := m.repo.GetUserByEmail(sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userWithData, err := m.repo.GetUserWithDataById(user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't get user data"})

		}

		userInfo := &dto.UserInfoDto{
			Name:            userWithData.UserData.Name,
			Birthdate:       userWithData.UserData.Birthdate,
			Email:           userWithData.UserData.Email,
			PhoneNumber:     userWithData.UserData.PhoneNumber,
			Telegram:        userWithData.UserData.Telegram,
			City:            userWithData.UserData.City,
			Age:             userWithData.UserData.Age,
			Employment:      userWithData.UserData.Employment,
			IsBusinessOwner: userWithData.UserData.IsBusinessOwner,
			PositionAtWork:  userWithData.UserData.PositionAtWork,
			MonthIncome:     userWithData.UserData.MonthIncome,
		}

		c.Set("user", userInfo)
		c.Set("self", user.ID)

		c.Next()
		return
	}
}

func (m *Middleware) AdminVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext, ok := c.Get("user")
		userDto := userContext.(*dto.UserInfoDto)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is not specified"})
			return
		}

		user, err := m.repo.GetUserByEmail(userDto.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		}

		if user.Role != int(service.Admin) || user.Role < 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No admin privilegies"})
		}

		c.Next()
	}
}

// func (m *Middleware) SelfVerificationMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userContext, ok := c.Get("user")
// 		userDto := userContext.(*dto.UserInfoDto)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is not specified"})
// 			return
// 		}

// 		user, err := m.repo.GetUserByEmail(userDto.Email)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			return
// 		}

// 		self, ok := c.Get("self")
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Self is not specified"})
// 			return
// 		}

// 		if self != user.ID {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No self privilegies"})
// 		}

// 		c.Next()
// 	}
// }

func (m *Middleware) CourseEnrollmentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		courseIDStr := c.Param("courseId")
		if courseIDStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
			return
		}

		courseID, err := uuid.Parse(courseIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
			return
		}

		userID, exists := c.Get("self")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		assignment, err := m.courseRepo.GetCourseAssignment(courseID, userID.(uuid.UUID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is not enrolled in this course"})
			return
		}

		c.Set("courseAssignment", assignment)

		c.Next()
	}
}
