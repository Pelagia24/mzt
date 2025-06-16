package repository

import (
	"mzt/config"
	"mzt/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// интерфейс для работы с платежами
// определяет все методы которые нужны для работы с платежами и ценами курсов в базе
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

// репозиторий для работы с платежами
// реализует интерфейс PaymentRepository
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

// CreatePayment создает новый платеж
// просто создает новую запись в таблице payments
func (r *PaymentRepo) CreatePayment(payment *entity.Payment) error {
	return r.DB.Create(payment).Error
}

// GetPaymentByID получает информацию о платеже
// ищет платеж в базе по его id
func (r *PaymentRepo) GetPaymentByID(paymentID uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.DB.Where("payment_id = ?", paymentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetPaymentsByUserID получает список платежей пользователя
// берет все платежи пользователя из базы
func (r *PaymentRepo) GetPaymentsByUserID(userID uuid.UUID) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	err := r.DB.Where("user_id = ?", userID).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// GetPaymentsByCourseID получает список платежей по курсу
// берет все платежи для конкретного курса из базы
func (r *PaymentRepo) GetPaymentsByCourseID(courseID uuid.UUID) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	err := r.DB.Where("course_id = ?", courseID).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// UpdatePaymentStatus обновляет статус платежа
// меняет статус платежа в базе на новый
func (r *PaymentRepo) UpdatePaymentStatus(paymentID uuid.UUID, status string) error {
	return r.DB.Model(&entity.Payment{}).Where("payment_id = ?", paymentID).Update("status", status).Error
}

// GetCoursePrice получает цену курса
// берет цену курса из базы по его id
func (r *PaymentRepo) GetCoursePrice(courseID uuid.UUID) (*entity.CoursePrice, error) {
	var price entity.CoursePrice
	err := r.DB.Where("course_id = ?", courseID).First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

// SetCoursePrice устанавливает цену курса
// создает новую запись с ценой курса в базе
func (r *PaymentRepo) SetCoursePrice(price *entity.CoursePrice) error {
	return r.DB.Create(price).Error
}

// UpdateCoursePrice обновляет цену курса
// меняет цену курса в базе на новую
func (r *PaymentRepo) UpdateCoursePrice(courseID uuid.UUID, amount float64) error {
	return r.DB.Model(&entity.CoursePrice{}).Where("course_id = ?", courseID).Update("amount", amount).Error
}
