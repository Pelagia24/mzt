package service

import (
	"errors"
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"mzt/internal/repository"
	"mzt/internal/validator"

	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	Default Role = iota // обычный пользователь
	Admin               // администратор
)

// преобразуем роль в строку
func (r Role) String() string {
	return [...]string{"Default", "Admin"}[r]
}

// сервис для работы с пользователями
type UserService struct {
	config    *config.Config
	repo      repository.UserRepository
	validator *validator.Validator
}

// создаем новый сервис для работы с пользователями(конструктор)
func NewUserService(cfg *config.Config, repo repository.UserRepository) *UserService {
	return &UserService{
		config:    cfg,
		repo:      repo,
		validator: validator.NewValidator(),
	}
}

// GetUserId получает id пользователя по почте
// просто ищет пользователя в базе и возвращает его id
func (s *UserService) GetUserId(email string) (uuid.UUID, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

// SignUp регистрирует нового пользователя
// создает нового пользователя в базе и выдает токены
func (s *UserService) SignUp(user *dto.RegistrationDto) (string, string, error) {
	// хешируем пароль чтобы не хранить его в открытом виде
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	// создаем нового пользователя
	userID := uuid.New()
	userEntity := entity.User{
		ID:         userID,
		PasswdHash: string(hashedPassword),
	}
	userData := entity.UserData{
		UserID:          userID,
		Email:           user.Email,
		Name:            user.Name,
		Birthdate:       user.Birthdate,
		PhoneNumber:     user.PhoneNumber,
		Telegram:        user.Telegram,
		City:            user.City,
		Age:             user.Age,
		Employment:      user.Employment,
		IsBusinessOwner: user.IsBusinessOwner,
		PositionAtWork:  user.PositionAtWork,
		MonthIncome:     user.MonthIncome,
	}

	// создаем токены для нового пользователя
	access, refresh, err := s.generateTokens(user.Email)
	if err != nil {
		return "", "", err
	}

	// сохраняем refresh токен в базе
	userAuth := entity.Auth{
		UserID: userID,
		Key:    refresh,
	}

	// сохраняем пользователя в базе
	err = s.repo.CreateUser(&userEntity, &userData, &userAuth)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// SignIn входит в аккаунт пользователя
// проверяет пароль и выдает новые токены
func (s *UserService) SignIn(user *dto.LoginDto) (string, string, error) {
	// ищем пользователя по почте
	userEntity, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return "", "", err
	}
	// проверяем что пароль правильный
	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.PasswdHash), []byte(user.Password)); err != nil {
		return "", "", err
	}

	// получаем данные пользователя
	userWithData, err := s.repo.GetUserWithDataById(userEntity.ID)
	if err != nil {
		return "", "", err
	}

	// создаем новые токены
	access, refresh, err := s.generateTokens(userWithData.UserData.Email)
	if err != nil {
		return "", "", err
	}

	// обновляем refresh токен в базе
	err = s.repo.UpdateToken(userEntity.ID, refresh)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// GetUsers получает список всех пользователей
// берет всех пользователей из базы и преобразует их в формат для ответа
func (s *UserService) GetUsers() ([]dto.UserInfoAdminDto, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserInfoAdminDto, 0)

	// для каждого пользователя собираем информацию о его курсах
	for _, user := range users {
		courseAssignments := make([]dto.CourseDto, 0)
		for _, assignment := range user.CourseAssignments {
			courseAssignments = append(courseAssignments, dto.CourseDto{
				CourseID: assignment.CourseID,
				Name:     assignment.Course.Title,
			})
		}

		// добавляем пользователя в ответ
		result = append(result, dto.UserInfoAdminDto{
			ID:                user.ID,
			Name:              user.UserData.Name,
			Birthdate:         user.UserData.Birthdate,
			Email:             user.UserData.Email,
			PhoneNumber:       user.UserData.PhoneNumber,
			Telegram:          user.UserData.Telegram,
			City:              user.UserData.City,
			Age:               user.UserData.Age,
			Employment:        user.UserData.Employment,
			IsBusinessOwner:   user.UserData.IsBusinessOwner,
			PositionAtWork:    user.UserData.PositionAtWork,
			MonthIncome:       user.UserData.MonthIncome,
			CourseAssignments: courseAssignments,
		})
	}
	return result, nil
}

// GetUser получает информацию о пользователе
// берет пользователя из базы и преобразует его в формат для ответа
func (s *UserService) GetUser(userId uuid.UUID) (*dto.UserInfoAdminDto, error) {
	user, err := s.repo.GetUserWithDataById(userId)
	if err != nil {
		return nil, err
	}
	userDto := &dto.UserInfoAdminDto{
		ID:                user.ID,
		Name:              user.UserData.Name,
		Birthdate:         user.UserData.Birthdate,
		Email:             user.UserData.Email,
		PhoneNumber:       user.UserData.PhoneNumber,
		Telegram:          user.UserData.Telegram,
		City:              user.UserData.City,
		Age:               user.UserData.Age,
		Employment:        user.UserData.Employment,
		IsBusinessOwner:   user.UserData.IsBusinessOwner,
		PositionAtWork:    user.UserData.PositionAtWork,
		MonthIncome:       user.UserData.MonthIncome,
		CourseAssignments: nil,
	}

	return userDto, nil
}

// UpdateUser обновляет информацию о пользователе
// меняет данные пользователя в базе
func (s *UserService) UpdateUser(userId uuid.UUID, updated *dto.UpdateUserDto) error {
	// считаем возраст пользователя
	age := time.Since(updated.Birthdate).Hours() / 24 / 365.25
	updated.Age = uint(age)

	// обновляем данные пользователя
	updatedEnity := &entity.UserData{
		UserID:          userId,
		Name:            updated.Name,
		Birthdate:       updated.Birthdate,
		Email:           updated.Email,
		PhoneNumber:     updated.PhoneNumber,
		Telegram:        updated.Telegram,
		City:            updated.City,
		Age:             updated.Age,
		Employment:      updated.Employment,
		IsBusinessOwner: updated.IsBusinessOwner,
		PositionAtWork:  updated.PositionAtWork,
		MonthIncome:     updated.MonthIncome,
	}
	return s.repo.UpdateUser(userId, updatedEnity)
}

// DeleteUser удаляет пользователя
// просто удаляет пользователя из базы
func (s *UserService) DeleteUser(toDel uuid.UUID) error {
	err := s.repo.DeleteUser(toDel)
	if err != nil {
		return err
	}
	return nil
}

// Role получает роль пользователя
// просто берет роль пользователя из базы
func (s *UserService) Role(userId uuid.UUID) (string, error) {
	user, err := s.repo.GetUserWithDataById(userId)
	if err != nil {
		return "", err
	}

	return Role(user.Role).String(), nil
}

/**
общие понятия:
- Access Token -- краткоживущий токен, используемый для аутентификации и авторизации пользователя
- Refresh Token -- долгоживущий токен, позволяющий получить новый Access Token без повторной авторизации
- cookie(как имя параметра метода) -- это строка, содержащая Refresh Token, переданная от клиента
*/

// RefreshTokens обновляет пару access/refresh токенов по переданному refresh token (cookie)
func (s *UserService) RefreshTokens(cookie string) (string, string, error) {
	// валидирует переданный токен с использованием refresh-ключа из токена
	token, err := s.validator.ValidateToken(cookie, s.config.Jwt.RefreshKey)
	if err != nil || !token.Valid {
		return "", "", err // если токен невалиден или возникла ошибка -- возвращается ошибка
	}

	// получает subject (обычно email) из токена
	sub, err := token.Claims.GetSubject()
	if err != nil || sub == "" {
		return "", "", err // ошибка получения subject или subject пустой
	}

	// получает пользователя по email (subject)
	user, err := s.repo.GetUserByEmail(sub)
	if err != nil {
		return "", "", err // ошибка получения пользователя из репозитория
	}

	// получает дополнительные данные(фио, telegram и т д) пользователя по его ID
	userData, err := s.repo.GetUserWithDataById(user.ID)
	if err != nil {
		return "", "", err // ошибка получения данных
	}

	// получает сущность пользователя с сохранённым refresh токеном (user.Auth.Key)
	userEntity, err := s.repo.GetUserWithRefreshById(user.ID)
	if err != nil {
		return "", "", err // ошибка получения записи из бд
	}

	// проверяет, совпадает ли переданный токен с сохранённым -- защита от повторного использования
	if userEntity.Auth.Key != cookie {
		return "", "", errors.New("this refresh token was already refreshed")
	}

	// генерирует новую пару access и refresh токенов
	access, refresh, err := s.generateTokens(userData.UserData.Email)
	if err != nil {
		return "", "", err // ошибка генерации токенов
	}

	// обновляет сохранённый refresh токен в базе
	if err := s.repo.UpdateToken(userEntity.ID, refresh); err != nil {
		return "", "", err // ошибка обновления токена в хранилище
	}

	// возвращает новую пару access и refresh токенов
	return access, refresh, nil
}

// generateTokens создает новые токены для пользователя
// создает access и refresh токены с указанной почтой
func (s *UserService) generateTokens(email string) (access string, refresh string, error error) {
	// создаем access токен
	access, err := s.validator.GenerateToken(email, s.config.Jwt.AccessKey, s.config.Jwt.AccessExpiresIn)
	if err != nil {
		return "", "", err
	}

	// создаем refresh токен
	refresh, err = s.validator.GenerateToken(email, s.config.Jwt.RefreshKey, s.config.Jwt.RefreshExpiresIn)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// Logout выходит из аккаунта
// просто удаляет refresh токен из базы
func (s *UserService) Logout(userId uuid.UUID) error {
	return s.repo.UpdateToken(userId, "")
}
