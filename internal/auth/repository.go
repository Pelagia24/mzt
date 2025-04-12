package auth

import (
	"fmt"
	"mzt/config"
	"mzt/internal/auth/entity"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepo struct {
	config *config.Config
	DB     *gorm.DB
}

type UserRepository interface {
	CreateUser()
	UpdateUser()
	GetRefreshTokenByEmail()
	// GetInternalIdByEmail()
	GetUserByEmail()
}

func NewRefreshTokensRepo(cfg *config.Config) *UserRepo {
	db := connectDB(cfg)
	return &UserRepo{
		config: cfg,
		DB:     db,
	}
}

func (r *UserRepo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(user *entity.User, userData *entity.UserData, auth *entity.Auth) error {

	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Create(&userData).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Create(&auth).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *UserRepo) GetUserWithDataByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("UserData").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserWithRefreshByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Auth").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateToken(userID uuid.UUID, token string) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var auth entity.Auth
	if err := tx.Where("user_id = ?", userID).First(&auth).Error; err != nil {
		tx.Rollback()
		return err
	}

	auth.Key = token

	if err := tx.Save(&auth).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func connectDB(config *config.Config) *gorm.DB {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Name,
		config.DB.Password)

	postrgresDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return postrgresDB
}
