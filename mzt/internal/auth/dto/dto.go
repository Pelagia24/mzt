package dto

import "time"

type User struct {
	Name            string    `json:"name" binding:"required"`
	Birthdate       time.Time `json:"birthdate" binding:"required"`
	Email           string    `json:"email" binding:"required"`
	PhoneNumber     string    `json:"phone_number" binding:"required"`
	Telegram        string    `json:"telegram" binding:"required"`
	City            string    `json:"city" binding:"required"`
	Age             uint      `json:"age" binding:"required"`
	Employment      string    `json:"employment" binding:"required"`
	IsBusinessOwner string    `json:"is_business_owner" binding:"required"`
	PositionAtWork  string    `json:"position_at_work" binding:"required"`
	MonthIncome     uint      `json:"month_income" binding:"required"`
	Password        string    `json:"password" binding:"required"`
	//ConfirmPassword string    `json:"confirm_password"`
}

type UserInfo struct {
	Name            string    `json:"name" binding:"required"`
	Birthdate       time.Time `json:"birthdate" binding:"required"`
	Email           string    `json:"email" binding:"required"`
	PhoneNumber     string    `json:"phone_number" binding:"required"`
	Telegram        string    `json:"telegram" binding:"required"`
	City            string    `json:"city" binding:"required"`
	Age             uint      `json:"age" binding:"required"`
	Employment      string    `json:"employment" binding:"required"`
	IsBusinessOwner string    `json:"is_business_owner" binding:"required"`
	PositionAtWork  string    `json:"position_at_work" binding:"required"`
	MonthIncome     uint      `json:"month_income" binding:"required"`
}

type UserLoginInfo struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
