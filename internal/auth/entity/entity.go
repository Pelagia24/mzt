package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserJWT struct {
	User   User      `gorm:"foreignKey:UserId"`
	UserId uuid.UUID `gorm:"type:uuid;not null"`
	Key    string    `gorm:"type:varchar(255);not null"`
}

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name            string    `gorm:"type:varchar(255);not null"`
	Birthdate       time.Time `gorm:"type:timestamp;not null"`
	Email           string    `gorm:"type:varchar(255);unique;not null"`
	PhoneNumber     string    `gorm:"type:varchar(255);not null"`
	Telegram        string    `gorm:"type:varchar(255)"`
	City            string    `gorm:"type:varchar(255);not null"`
	Age             uint      `gorm:"not null"`
	Employment      string    `gorm:"type:varchar(255);not null"`
	IsBusinessOwner string    `gorm:"type:varchar(255);not null"`
	PositionAtWork  string    `gorm:"type:varchar(255);not null"`
	MonthIncome     uint      `gorm:"not null"`
	PasswordHash    string    `gorm:"type:varchar(255);not null"`
}
