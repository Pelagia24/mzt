package service

import (
	"errors"
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"mzt/internal/repository"
	"mzt/internal/validator"

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

type UserService struct {
	config    *config.Config
	repo      repository.UserRepository
	validator *validator.Validator
}

func NewUserService(cfg *config.Config, repo repository.UserRepository) *UserService {
	return &UserService{
		config:    cfg,
		repo:      repo,
		validator: validator.NewValidator(),
	}
}

func (s *UserService) GetUserId(email string) (uuid.UUID, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (s *UserService) SignUp(user *dto.RegistrationDto) (string, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

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

	access, refresh, err := s.generateTokens(user.Email)
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

func (s *UserService) SignIn(user *dto.LoginDto) (string, string, error) {
	userEntity, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.PasswdHash), []byte(user.Password)); err != nil {
		return "", "", err
	}

	userWithData, err := s.repo.GetUserWithDataById(userEntity.ID)
	if err != nil {
		return "", "", err
	}

	access, refresh, err := s.generateTokens(userWithData.UserData.Email)
	if err != nil {
		return "", "", err
	}

	err = s.repo.UpdateToken(userEntity.ID, refresh)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *UserService) GetUsers() ([]dto.UserInfoAdminDto, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserInfoAdminDto, 0)

	for _, user := range users {
		courseAssignments := make([]dto.CourseDto, 0)
		for _, assignment := range user.CourseAssignments {
			courseAssignments = append(courseAssignments, dto.CourseDto{
				CourseID: assignment.CourseID,
				Name:     assignment.Course.Title,
			})
		}

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

func (s *UserService) UpdateUser(userId uuid.UUID, updated *dto.UpdateUserDto) error {
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

func (s *UserService) DeleteUser(toDel uuid.UUID) error {
	err := s.repo.DeleteUser(toDel)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Role(userId uuid.UUID) (string, error) {
	user, err := s.repo.GetUserWithDataById(userId)
	if err != nil {
		return "", err
	}

	return Role(user.Role).String(), nil
}

func (s *UserService) RefreshTokens(cookie string) (string, string, error) {
	token, err := s.validator.ValidateToken(cookie, s.config.Jwt.RefreshKey)
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

	userData, err := s.repo.GetUserWithDataById(user.ID)
	if err != nil {
		return "", "", err
	}

	userEntity, err := s.repo.GetUserWithRefreshById(user.ID)
	if err != nil {
		return "", "", err
	}

	if userEntity.Auth.Key != cookie {
		return "", "", errors.New("this refresh token was already refreshed")
	}

	access, refresh, err := s.generateTokens(userData.UserData.Email)
	if err != nil {
		return "", "", err
	}

	if err := s.repo.UpdateToken(userEntity.ID, refresh); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *UserService) generateTokens(email string) (access string, refresh string, error error) {
	access, err := s.validator.GenerateToken(email, s.config.Jwt.AccessKey, s.config.Jwt.AccessExpiresIn)
	if err != nil {
		return "", "", err
	}

	refresh, err = s.validator.GenerateToken(email, s.config.Jwt.RefreshKey, s.config.Jwt.RefreshExpiresIn)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
