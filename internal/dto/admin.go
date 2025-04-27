package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserInfoAdminDto struct {
	ID                uuid.UUID   `json:"id" binding:"required"`
	Name              string      `json:"name" binding:"required"`
	Birthdate         time.Time   `json:"birthdate" binding:"required"`
	Email             string      `json:"email" binding:"required"`
	PhoneNumber       string      `json:"phone_number" binding:"required"`
	Telegram          string      `json:"telegram" binding:"required"`
	City              string      `json:"city" binding:"required"`
	Age               uint        `json:"age"`
	Employment        string      `json:"employment" binding:"required"`
	IsBusinessOwner   string      `json:"is_business_owner" binding:"required"`
	PositionAtWork    string      `json:"position_at_work" binding:"required"`
	MonthIncome       uint        `json:"month_income"`
	CourseAssignments []CourseDto `json:"course_assignments"`
}
type CreateCourseDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       uint   `json:"price" binding:"required"`
}

type UpdateCourseDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       uint   `json:"price" binding:"required"`
}

type CreateLessonDto struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	VideoURL    string `json:"video_url" binding:"required"`
	SummaryURL  string `json:"summary_url" binding:"required"`
}

type UpdateLessonDto struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	VideoURL    string `json:"video_url" binding:"required"`
	SummaryURL  string `json:"summary_url" binding:"required"`
}
type UpdateUserDto struct {
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
}
