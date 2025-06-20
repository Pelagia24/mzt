package router

import (
	"net/http"

	"mzt/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получает список всех пользователей
// доступно только админам
func (r *Router) GetUsers(c *gin.Context) {
	// получаем список пользователей из сервиса
	users, err := r.authService.GetUsers()
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't get users"})
		return
	}

	// отправляем список пользователей клиенту
	c.JSON(http.StatusOK, gin.H{
		"message": "User infos",
		"users":   users,
	})
}

// получает роль пользователя по его id
// доступно только админам
func (r *Router) Role(c *gin.Context) {
	// достаем id пользователя из параметров запроса
	userId := c.Param("user_id")
	id, err := uuid.Parse(userId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// получаем роль пользователя из сервиса
	role, err := r.authService.Role(id)
	if err != nil || role == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't get user role"})
		return
	}
	// отправляем роль пользователя клиенту
	c.JSON(http.StatusOK, gin.H{
		"message": "Role of user",
		"role":    role,
	})
}

// обрабатывает запросы на получение, обновление и удаление пользователя
// доступно только админам
func (r *Router) Users(c *gin.Context) {
	// достаем id пользователя из параметров запроса
	userId := c.Param("user_id")
	id, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// в зависимости от метода запроса выполняем разные действия
	switch c.Request.Method {
	case http.MethodGet:
		// получаем информацию о пользователе
		user, err := r.authService.GetUser(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Can't get user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User detailedinfo",
			"user":    user,
		})

	case http.MethodDelete:
		// достаем id текущего пользователя из контекста
		self, ok := c.Get("self")
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unknown sender"})
			return
		}

		selfId, ok := self.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unknown sender"})
			return
		}

		// проверяем что пользователь не пытается удалить сам себя
		if id == selfId {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't delete self"})
			return
		}

		// удаляем пользователя через сервис
		err = r.authService.DeleteUser(id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deleted user with success"})
	case http.MethodPut:
		// парсим данные из тела запроса
		var payload dto.UpdateUserDto
		if err := c.ShouldBindJSON(&payload); err != nil {
			// если данные невалидные, возвращаем ошибку
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// обновляем пользователя через сервис
		err := r.authService.UpdateUser(id, &payload)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Can't update user"})
			return
		}

		// отправляем успешный response
		c.JSON(http.StatusOK, gin.H{"message": "Updated user with success"})

	default:
		// если метод запроса не поддерживается, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid http method"})
		return
	}

}

// получает профиль текущего пользователя
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

// SignIn обрабатывает вход пользователя
// проверяет почту и пароль и выдает токены если все ок
func (r *Router) SignIn(c *gin.Context) {
	// берем данные из запроса
	var payload dto.LoginDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем что почта правильная
	if !r.validator.IsValidEmail(payload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// проверяем что пароль достаточно сложный
	if !r.validator.IsValidPassword(payload.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is weak or contains incorrect symbols"})
		return
	}

	// пытаемся войти и получить токены
	access, refresh, err := r.authService.SignIn(&payload)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// получаем id пользователя по почте
	id, err := r.authService.GetUserId(payload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// узнаем роль пользователя админ он или нет
	role, err := r.authService.Role(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// сохраняем refresh токен в куки
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", r.config.Jwt.Domain, false, true)

	// отправляем ответ с токенами и данными пользователя
	c.JSON(http.StatusOK, gin.H{
		"message":      "User signed in successfully",
		"access_token": access,
		"id":           id,
		"role":         role,
	})
}

// регистрирует нового пользователя
// проверяет все поля и создает аккаунт если все ок
func (r *Router) SignUp(c *gin.Context) {
	// берем данные из запроса
	var payload dto.RegistrationDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем что почта правильная
	if !r.validator.IsValidEmail(payload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// проверяем что пароль достаточно сложный
	if !r.validator.IsValidPassword(payload.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	// проверяем что имя нормальное
	if !r.validator.IsValidName(payload.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}

	// проверяем что телефон правильный
	if !r.validator.IsValidPhoneNumber(payload.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}

	// проверяем что телеграм правильный
	if !r.validator.IsValidTelegram(payload.Telegram) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telegram"})
		return
	}

	// создаем пользователя и получаем токены
	access, refresh, err := r.authService.SignUp(&payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// получаем id нового пользователя
	id, err := r.authService.GetUserId(payload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// узнаем роль пользователя админ он или нет
	role, err := r.authService.Role(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// сохраняем refresh токен в куки
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", r.config.Jwt.Domain, false, true)

	// отправляем ответ с токенами и данными пользователя
	c.JSON(http.StatusCreated, gin.H{
		"message":      "User created successfully",
		"access_token": access,
		"id":           id,
		"role":         role,
	})
}

// Refresh обновляет токены пользователя
// берет refresh токен из куки и выдает новые токены
func (r *Router) Refresh(c *gin.Context) {
	// достаем refresh токен из куки
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// получаем новые токены
	access, refresh, err := r.authService.RefreshTokens(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// проверяем что токен валидный
	parsed, err := r.validator.ValidateToken(token, r.config.Jwt.RefreshKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// достаем почту из токена
	sub, err := parsed.Claims.GetSubject()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// получаем id пользователя по почте
	id, err := r.authService.GetUserId(sub)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// узнаем роль пользователя админ он или нет
	role, err := r.authService.Role(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// сохраняем новый refresh токен в куки
	c.SetCookie("refresh_token", refresh, int(r.config.Jwt.RefreshExpiresIn.Seconds()), "/", r.config.Jwt.Domain, false, true)

	// отправляем ответ с новыми токенами и данными пользователя
	c.JSON(http.StatusOK, gin.H{
		"message":      "Tokens refreshed successfully",
		"access_token": access,
		"id":           id,
		"role":         role,
	})
}

// Logout выходит из аккаунта
// удаляет refresh токен и очищает куки
func (r *Router) Logout(c *gin.Context) {
	// достаем id пользователя из контекста
	user, ok := c.Get("self")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// удаляем refresh токен из базы
	userId := user.(uuid.UUID)
	err := r.authService.Logout(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// очищаем куки с refresh токеном
	c.SetCookie("refresh_token", "", -1, "/", r.config.Jwt.Domain, false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
