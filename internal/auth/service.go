package auth

import (
	"errors"
	"mzt/config"
	"mzt/internal/auth/dto"
	"mzt/internal/auth/entity"
	"mzt/internal/auth/utils"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	config *config.Config
	repo   *UserRepo
}

func NewService(cfg *config.Config, repo *UserRepo) *Service {
	return &Service{
		config: cfg,
		repo:   repo,
	}
}

func (s *Service) SignUp(user *dto.RegistrationDto) (string, string, error) {
	if !isValidEmail(user.Email) {
		return "", "", errors.New("Invalid Email")
	}

	if !isValidPhoneNumber(user.PhoneNumber) {
		return "", "", errors.New("Invalid Phone Number")
	}

	if !isValidTelegram(user.Telegram) {
		return "", "", errors.New("Invalid Telegram")
	}

	if user.Password != user.ConfirmPassword {
		return "", "", errors.New("Password and confirmation don't match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	userID := uuid.New()
	userEntity := entity.User{
		ID:    userID,
		Email: user.Email,
		// Role:
		PasswdHash: string(hashedPassword),
	}
	userData := entity.UserData{
		UserID:          userID,
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

	access, refresh, err := s.generateTokens(userEntity.Email)
	if err != nil {
		return "", "", err
	}

	userAuth := entity.Auth{
		UserID: userID,
		Key:    refresh,
	}

	err = s.repo.CreateUser(&userEntity, &userData, &userAuth)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
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

func (s *Service) SignIn(user *dto.LoginDto) (string, string, error) {
	if !isValidEmail(user.Email) {
		return "", "", errors.New("Invalid Email")
	}

	userEntity, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.PasswdHash), []byte(user.Password)); err != nil {
		return "", "", err
	}

	access, refresh, err := s.generateTokens(userEntity.Email)
	if err != nil {
		return "", "", err
	}

	err = s.repo.UpdateToken(userEntity.ID, refresh)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) RefreshTokens(cookie string) (string, string, error) {
	token, err := utils.ValidateToken(cookie, s.config.Jwt.RefreshKey)
	if err != nil || !token.Valid {
		return "", "", err
	}

	sub, err := token.Claims.GetSubject()
	if err != nil || sub == "" {
		return "", "", err
	}

	userEntity, err := s.repo.GetUserWithRefreshByEmail(sub)
	if err != nil {
		return "", "", err
	}

	if userEntity.Auth.Key != cookie {
		return "", "", errors.New("This refresh token already refreshed")
	}

	access, refresh, err := s.generateTokens(userEntity.Email)
	if err != nil {
		return "", "", err
	}

	if err := s.repo.UpdateToken(userEntity.ID, refresh); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) generateTokens(email string) (access string, refresh string, error error) {

	access, err := utils.GenerateToken(email, s.config.Jwt.AccessKey, s.config.Jwt.AccessExpiresIn)
	if err != nil {
		return "", "", err
	}

	refresh, err = utils.GenerateToken(email, s.config.Jwt.RefreshKey, s.config.Jwt.RefreshExpiresIn)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
