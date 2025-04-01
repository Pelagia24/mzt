package auth

import (
	"mzt/config"
	"mzt/internal/auth/entity"
	"mzt/internal/auth/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokensRepo struct {
	config *config.Config
	DB     *gorm.DB
}

func NewRefreshTokensRepo(config *config.Config) *RefreshTokensRepo {
	db := utils.ConnectDB(config)
	return &RefreshTokensRepo{
		config: config,
		DB:     db,
	}
}

func (r *RefreshTokensRepo) GetInternalIdByEmail(email string) (uuid.UUID, error) {
	var user entity.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (r *RefreshTokensRepo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RefreshTokensRepo) CreateUser(user *entity.User) error {
	err := r.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokensRepo) CreateUserJWT(userJWTEntity *entity.UserJWT) error {
	err := r.DB.Create(&userJWTEntity).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokensRepo) UpdateToken(userId uuid.UUID, token string) error {
	tx := r.DB.Begin()
	err := tx.Model(&entity.UserJWT{}).Where("user_id = ?", userId).Updates(entity.UserJWT{Key: token}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

// func (r *RefreshTokensRepo) RemoveToken(token string) error {
// 	var user entity.UserJWT
// 	err := r.DB.First(&user, "key = ?", token).Error
// 	if err != nil {
// 		return err
// 	}
// 	err = r.DB.Delete(&user).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *RefreshTokensRepo) GetUserJWTByToken(token string) (*entity.UserJWT, error) {
	var userJWT entity.UserJWT
	err := r.DB.Where("key = ?", token).First(&userJWT).Error
	if err != nil {
		return nil, err
	}
	return &userJWT, nil
}

func (r *RefreshTokensRepo) GetToken(userId string) (string, error) {
	var user *entity.UserJWT
	err := r.DB.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Key, nil
}
