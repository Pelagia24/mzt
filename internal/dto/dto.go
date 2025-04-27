package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegistrationDto struct {
	Name        string    `json:"name" binding:"required"`
	Birthdate   time.Time `json:"birthdate" binding:"required"`
	Email       string    `json:"email" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Telegram    string    `json:"telegram" binding:"required"`
	City        string    `json:"city" binding:"required"`
	Age         uint      `json:"age"`
	// Age             uint      `json:"age" binding:"required"`
	Employment      string `json:"employment" binding:"required"`
	IsBusinessOwner string `json:"is_business_owner" binding:"required"`
	PositionAtWork  string `json:"position_at_work" binding:"required"`
	MonthIncome     uint   `json:"month_income"`
	// MonthIncome     uint      `json:"month_income" binding:"required"`
	Password string `json:"password" binding:"required"`
	// ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type UserInfoDto struct {
	Name        string    `json:"name" binding:"required"`
	Birthdate   time.Time `json:"birthdate" binding:"required"`
	Email       string    `json:"email" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Telegram    string    `json:"telegram" binding:"required"`
	City        string    `json:"city" binding:"required"`
	Age         uint      `json:"age"`
	// Age             uint      `json:"age" binding:"required"`
	Employment      string `json:"employment" binding:"required"`
	IsBusinessOwner string `json:"is_business_owner" binding:"required"`
	PositionAtWork  string `json:"position_at_work" binding:"required"`
	MonthIncome     uint   `json:"month_income"`
	// MonthIncome     uint   `json:"month_income" binding:"required"`
}

type LoginDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LessonDto struct {
	LessonID uuid.UUID `json:"lesson_id"`
	CourseID uuid.UUID `json:"course_id"`
	Title    string    `json:"title"`
	Summery  string    `json:"summery"`
	VideoURL string    `json:"video_url"`
	Text     string    `json:"text"`
}

type AssignUserToCourseDto struct {
	UserId string `json:"user_id" binding:"required"`
}

type UpdateProgressDto struct {
	Progress uint `json:"progress" binding:"required"`
}

type CourseDto struct {
	CourseID    uuid.UUID `json:"course_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
