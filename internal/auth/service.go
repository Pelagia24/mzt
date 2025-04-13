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

type Role int

const (
	Default Role = iota
	Admin
)

func (r Role) String() string {
	return [...]string{"Default", "Admin"}[r]
}

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

	//TODO more validation

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

func (s *Service) GetUsers() ([]dto.UserInfoAdminDto, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserInfoAdminDto, 0)

	for _, user := range users {
		result = append(result, dto.UserInfoAdminDto{
			ID:                user.ID,
			Name:              user.UserData.Name,
			Birthdate:         user.UserData.Birthdate,
			Email:             user.Email,
			PhoneNumber:       user.UserData.PhoneNumber,
			Telegram:          user.UserData.Telegram,
			City:              user.UserData.City,
			Age:               user.UserData.Age,
			Employment:        user.UserData.Employment,
			IsBusinessOwner:   user.UserData.IsBusinessOwner,
			PositionAtWork:    user.UserData.PositionAtWork,
			MonthIncome:       user.UserData.MonthIncome,
			CourseAssignments: nil,
		})
	}
	//TODO check this out
	return result, nil
}

func (s *Service) GetUser(userId uuid.UUID) (*dto.UserInfoAdminDto, error) {
	user, err := s.repo.GetUserWithDataById(userId)
	if err != nil {
		return nil, err
	}
	userDto := &dto.UserInfoAdminDto{
		ID:                user.ID,
		Name:              user.UserData.Name,
		Birthdate:         user.UserData.Birthdate,
		Email:             user.Email,
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

func (s *Service) UpdateUser(userId uuid.UUID, updated *dto.UpdateUserDto) error {
	updatedEnity := &entity.UserData{
		UserID:          userId,
		Name:            updated.Name,
		Birthdate:       updated.Birthdate,
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

func (s *Service) DeleteUser(toDel uuid.UUID) error {
	err := s.repo.DeleteUser(toDel)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Role(userId uuid.UUID) (string, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return "", err
	}

	return Role(user.Role).String(), nil
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

	user, err := s.repo.GetUserByEmail(sub)
	if err != nil {
		return "", "", err
	}

	userEntity, err := s.repo.GetUserWithRefreshById(user.ID)
	if err != nil {
		return "", "", err
	}

	if userEntity.Auth.Key != cookie {
		return "", "", errors.New("This refresh token was already refreshed")
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
