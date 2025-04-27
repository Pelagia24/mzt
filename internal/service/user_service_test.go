package service

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_SignUp(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:        "test-access-key",
			RefreshKey:       "test-refresh-key",
			AccessExpiresIn:  time.Minute * 30,
			RefreshExpiresIn: time.Hour * 24 * 14,
		},
	}
	service := NewUserService(cfg, mockRepo)
	birthdate, _ := time.Parse("2006-01-02", "1990-01-01")
	userDto := &dto.RegistrationDto{
		Email:           "test@example.com",
		Password:        "password123",
		Name:            "Test User",
		Birthdate:       birthdate,
		PhoneNumber:     "+1234567890",
		Telegram:        "@testuser",
		City:            "Test City",
		Age:             30,
		Employment:      "Test Employment",
		IsBusinessOwner: "false",
		PositionAtWork:  "Test Position",
		MonthIncome:     5000,
	}

	access, refresh, err := service.SignUp(userDto)

	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	user, err := mockRepo.GetUserByEmail(userDto.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestService_GetUserId(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	cfg := &config.Config{}
	service := NewUserService(cfg, mockRepo)
	email := "test@example.com"

	birthdate, _ := time.Parse("2006-01-02", "1990-01-01")
	userDto := &dto.RegistrationDto{
		Email:           email,
		Password:        "password123",
		Name:            "Test User",
		Birthdate:       birthdate,
		PhoneNumber:     "+1234567890",
		Telegram:        "@testuser",
		City:            "Test City",
		Age:             30,
		Employment:      "Test Employment",
		IsBusinessOwner: "false",
		PositionAtWork:  "Test Position",
		MonthIncome:     5000,
	}
	_, _, err := service.SignUp(userDto)
	assert.NoError(t, err)

	userId, err := service.GetUserId(email)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, userId)
}

func TestService_RefreshTokens(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:        "test-access-key",
			RefreshKey:       "test-refresh-key",
			AccessExpiresIn:  time.Minute * 30,
			RefreshExpiresIn: time.Hour * 24 * 14,
		},
	}
	service := NewUserService(cfg, mockRepo)
	birthdate, _ := time.Parse("2006-01-02", "1990-01-01")
	userDto := &dto.RegistrationDto{
		Email:           "test@example.com",
		Password:        "password123",
		Name:            "Test User",
		Birthdate:       birthdate,
		PhoneNumber:     "+1234567890",
		Telegram:        "@testuser",
		City:            "Test City",
		Age:             30,
		Employment:      "Test Employment",
		IsBusinessOwner: "false",
		PositionAtWork:  "Test Position",
		MonthIncome:     5000,
	}

	_, refresh, err := service.SignUp(userDto)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
	newAccess, newRefresh, err := service.RefreshTokens(refresh)

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccess)
	assert.NotEmpty(t, newRefresh)
	assert.NotEqual(t, refresh, newRefresh)
}
