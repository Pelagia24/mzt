package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_userid_userdata"`
	Email           string    `gorm:"uniqueIndex:idx_email;not null"`
	Name            string
	Birthdate       time.Time
	PhoneNumber     string
	Telegram        string
	City            string
	Age             uint
	Employment      string
	IsBusinessOwner string
	PositionAtWork  string
	MonthIncome     uint
}

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	PasswdHash string
	//TODO make role an entity
	Role              int
	Auth              *Auth              `gorm:"constraint:OnDelete:CASCADE;"`
	UserData          *UserData          `gorm:"constraint:OnDelete:CASCADE;"`
	CourseAssignments []CourseAssignment `gorm:"constraint:OnDelete:CASCADE;"`
	// EventRecords []EventAssignment
}

type Auth struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_userid_auth"`
	Key    string    `gorm:"type:varchar(255);not null"`
}

type Course struct {
	CourseID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title    string
	Desc     string

	Lessons           []Lesson           `gorm:"constraint:OnDelete:CASCADE;"`
	CourseAssignments []CourseAssignment `gorm:"constraint:OnDelete:CASCADE;"`
}

type CourseAssignment struct {
	CaID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	CourseID uuid.UUID `gorm:"type:uuid;not null"`
	Progress uint

	User   User
	Course Course
}

type Lesson struct {
	LessonID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CourseID uuid.UUID `gorm:"type:uuid;not null"`
	Title    string
	Summery  string
	VideoURL string
	Text     string

	Course Course
}
