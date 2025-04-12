package auth

import (
	"mzt/config"
	"mzt/internal/auth/dto"
	"mzt/internal/auth/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	config *config.Config
	repo   *UserRepo
}

func NewMiddleware(config *config.Config, repo *UserRepo) *Middleware {
	return &Middleware{
		config: config,
		repo:   repo,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(AuthHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		tokenString := fields[1]
		token, err := utils.ValidateToken(tokenString, m.config.Jwt.AccessKey)

		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil || sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		//TODO: Implement dto that contains user information
		user, err := m.repo.GetUserWithDataByEmail(sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userInfo := &dto.UserInfoDto{
			Name:            user.UserData.Name,
			Birthdate:       user.UserData.Birthdate,
			Email:           user.Email,
			PhoneNumber:     user.UserData.PhoneNumber,
			Telegram:        user.UserData.Telegram,
			City:            user.UserData.City,
			Age:             user.UserData.Age,
			Employment:      user.UserData.Employment,
			IsBusinessOwner: user.UserData.IsBusinessOwner,
			PositionAtWork:  user.UserData.PositionAtWork,
			MonthIncome:     user.UserData.MonthIncome,
		}

		c.Set("user", userInfo)

		c.Next()
		return
	}
}
