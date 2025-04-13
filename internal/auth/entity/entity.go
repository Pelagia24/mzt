package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_userid_userdata"`
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
	Email      string    `gorm:"uniqueIndex:idx_email;not null"`
	PasswdHash string
	//TODO make role an entity
	Role              int
	Auth              *Auth
	UserData          *UserData
	CourseAssignments []CourseAssignment
	// EventRecords []EventAssignment
}

type Auth struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_userid_auth"`
	Key    string    `gorm:"type:varchar(255);not null"`
}

// Courses
type Course struct {
	CourseID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title    string
	Desc     string

	Lessons           []Lesson
	CourseAssignments []CourseAssignment
}

type CourseAssignment struct {
	CaID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	CourseID uuid.UUID `gorm:"type:uuid;not null"`
	Progress uint

	User   User   `gorm:"constraint:OnDelete:CASCADE;"`
	Course Course `gorm:"constraint:OnDelete:CASCADE;"`
}

type Lesson struct {
	LessonID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	CourseID   uuid.UUID `gorm:"type:uuid;not null"`
	VideoURL   string
	SummaryURL string

	Course Course `gorm:"constraint:OnDelete:CASCADE;"`
}
