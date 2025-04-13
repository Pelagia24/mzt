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
	CourseAssignments interface{} `json:"course_assignments"`
}
