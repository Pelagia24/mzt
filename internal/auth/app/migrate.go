package app

import (
	"errors"
	"fmt"
	"mzt/internal/auth"
	"mzt/internal/auth/entity"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate(r *auth.UserRepo) {
	r.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	err := r.DB.AutoMigrate(
		&entity.User{},
		&entity.UserData{},
		&entity.Auth{},
		&entity.Course{},
		&entity.CourseAssignment{},
		&entity.Lesson{},
	)
	if err != nil {
		panic(err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	userId := uuid.New()

	userEntity := &entity.User{
		ID:         userId,
		Role:       1,
		PasswdHash: string(hashedPassword),
	}

	userData := &entity.UserData{
		ID:          0,
		UserID:      userId,
		Email:       "test@test.test",
		Name:        "test",
		Birthdate:   time.Now(),
		PhoneNumber: "+71111111111",
		Telegram:    "@test123",
	}

	auth := &entity.Auth{
		ID:     0,
		UserID: userId,
		Key:    "",
	}

	var user entity.User
	err = r.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = r.CreateUser(userEntity, userData, auth)
			if err != nil {
				panic(err)
			}
			fmt.Println("migrated")
		} else {
			panic(err)
		}

	}
}
