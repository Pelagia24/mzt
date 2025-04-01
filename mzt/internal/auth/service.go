package auth

import (
	"mzt/config"
	"mzt/internal/auth/dto"
	"mzt/internal/auth/entity"
	"mzt/internal/auth/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	config     *config.Config
	repository *RefreshTokensRepo
}

func NewService(config *config.Config, repository *RefreshTokensRepo) *Service {
	return &Service{
		config:     config,
		repository: repository,
	}
}

func (s *Service) SignUp(user *dto.User) (string, string, error) {
	//TODO: validation

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	userEntity := entity.User{
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
		PasswordHash:    string(hashedPassword),
	}

	err = s.repository.CreateUser(&userEntity)
	if err != nil {
		return "", "", err
	}

	id, err := s.repository.GetInternalIdByEmail(user.Email)
	if err != nil {
		return "", "", err
	}

	access, refresh, err := s.generateTokens(userEntity.Email)
	if err != nil {
		return "", "", err
	}

	userJWTEntity := &entity.UserJWT{
		User:   userEntity,
		UserId: id,
		Key:    refresh,
	}

	err = s.repository.CreateUserJWT(userJWTEntity)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) SignIn(userInfo *dto.UserLoginInfo) (string, string, error) {
	userEntity, err := s.repository.GetUserByEmail(userInfo.Email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.PasswordHash), []byte(userInfo.Password)); err != nil {
		return "", "", err
	}

	access, refresh, err := s.generateTokens(userEntity.Email)
	if err != nil {
		return "", "", err
	}

	userJWTEntity := &entity.UserJWT{
		User:   *userEntity,
		UserId: userEntity.ID,
		Key:    refresh,
	}

	err = s.repository.UpdateToken(userJWTEntity.UserId, userJWTEntity.Key)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) RefreshTokens(cookie string) (string, string, error) {
	_, err := utils.ValidateToken(cookie, s.config.Jwt.RefreshKey)
	if err != nil {
		return "", "", err
	}

	userJWTEntity, err := s.repository.GetUserJWTByToken(cookie)
	if err != nil {
		return "", "", err
	}
	//TODO better validation token sub and db entry
	userEntity, err := s.repository.GetUserById(userJWTEntity.UserId)
	if err != nil {
		return "", "", err
	}

	access, refresh, err := s.generateTokens(userEntity.Email)
	if err != nil {
		return "", "", err
	}

	if err := s.repository.UpdateToken(userJWTEntity.UserId, refresh); err != nil {
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
