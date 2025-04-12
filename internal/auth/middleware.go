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
	repo   *RefreshTokensRepo
}

func NewMiddleware(config *config.Config, repo *RefreshTokensRepo) *Middleware {
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
		user, err := m.repo.GetUserByEmail(sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userInfo := &dto.UserInfoDto{
			Name:            user.Name,
			Birthdate:       user.Birthdate,
			Email:           user.Email,
			PhoneNumber:     user.PhoneNumber,
			Telegram:        user.Telegram,
			City:            user.City,
			Age:             user.Age,
			Employment:      user.Employment,
			IsBusinessOwner: user.IsBusinessOwner,
			PositionAtWork:  user.PositionAtWork,
			MonthIncome:     user.MonthIncome,
		}

		c.Set("user", userInfo)

		c.Next()
		return
	}
}
