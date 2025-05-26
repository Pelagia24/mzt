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

	Lessons  []Lesson           `gorm:"constraint:OnDelete:CASCADE;"`
	Users    []CourseAssignment `gorm:"constraint:OnDelete:CASCADE;"`
	Events   []Event            `gorm:"constraint:OnDelete:CASCADE;"`
	Payments []Payment          `gorm:"constraint:OnDelete:CASCADE;"`
	Price    *CoursePrice       `gorm:"constraint:OnDelete:CASCADE;"`
}

// TODO index on entries, refund if error
// TODO also create payment repository
type CourseAssignment struct {
	CaID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index:,unique,composite:idx_user_course"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;index:,unique,composite:idx_user_course"`
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

type Event struct {
	EventID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	CourseID    uuid.UUID `gorm:"type:uuid;not null;index:idx_course_event"`
	Title       string    `gorm:"not null"`
	Description string
	EventDate   time.Time `gorm:"index:idx_event_date;not null"`
	SecretInfo  string

	Course Course `gorm:"constraint:OnDelete:CASCADE;"`
}

type Payment struct {
	PaymentID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index:idx_user_payment"`
	CourseID     uuid.UUID `gorm:"type:uuid;not null;index:idx_course_payment"`
	Amount       float64   `gorm:"not null"`
	CurrencyCode string    `gorm:"not null;default:'RUB'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Status       string    `gorm:"not null;default:'pending'"`
	PaymentRef   string    `gorm:"type:varchar(255)"`

	User   User   
	Course Course
}

type CoursePrice struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CourseID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_course_price"`
	Amount       float64   `gorm:"not null"`
	CurrencyCode string    `gorm:"not null;default:'RUB'"`

	Course Course 
}
