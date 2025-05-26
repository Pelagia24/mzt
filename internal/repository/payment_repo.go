package repository

import (
	"mzt/config"
	"mzt/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment *entity.Payment) error
	GetPaymentByID(paymentID uuid.UUID) (*entity.Payment, error)
	GetPaymentsByUserID(userID uuid.UUID) ([]*entity.Payment, error)
	GetPaymentsByCourseID(courseID uuid.UUID) ([]*entity.Payment, error)
	UpdatePaymentStatus(paymentID uuid.UUID, status string) error

	GetCoursePrice(courseID uuid.UUID) (*entity.CoursePrice, error)
	SetCoursePrice(price *entity.CoursePrice) error
	UpdateCoursePrice(courseID uuid.UUID, amount float64) error
}

type PaymentRepo struct {
	config *config.Config
	DB     *gorm.DB
}

func NewPaymentRepo(cfg *config.Config) *PaymentRepo {
	return &PaymentRepo{
		config: cfg,
		DB:     connectDB(cfg),
	}
}

func (r *PaymentRepo) CreatePayment(payment *entity.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *PaymentRepo) GetPaymentByID(paymentID uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.DB.Where("payment_id = ?", paymentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepo) GetPaymentsByUserID(userID uuid.UUID) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	err := r.DB.Where("user_id = ?", userID).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepo) GetPaymentsByCourseID(courseID uuid.UUID) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	err := r.DB.Where("course_id = ?", courseID).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepo) UpdatePaymentStatus(paymentID uuid.UUID, status string) error {
	return r.DB.Model(&entity.Payment{}).Where("payment_id = ?", paymentID).Update("status", status).Error
}

// Course price methods
func (r *PaymentRepo) GetCoursePrice(courseID uuid.UUID) (*entity.CoursePrice, error) {
	var price entity.CoursePrice
	err := r.DB.Where("course_id = ?", courseID).First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PaymentRepo) SetCoursePrice(price *entity.CoursePrice) error {
	return r.DB.Create(price).Error
}

func (r *PaymentRepo) UpdateCoursePrice(courseID uuid.UUID, amount float64) error {
	return r.DB.Model(&entity.CoursePrice{}).Where("course_id = ?", courseID).Update("amount", amount).Error
}
