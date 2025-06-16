// пакет для работы с сервисами
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"mzt/internal/repository"
	"net/http"

	"github.com/google/uuid"
)

// сервис для работы с платежами
// отвечает за создание платежей и работу с YooKassa
type PaymentService struct {
	config      *config.Config
	courseRepo  repository.CourseRepository
	paymentRepo repository.PaymentRepository
}

// создаем новый сервис для работы с платежами
func NewPaymentService(cfg *config.Config, courseRepo repository.CourseRepository, paymentRepo repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		config:      cfg,
		courseRepo:  courseRepo,
		paymentRepo: paymentRepo,
	}
}

// создает платеж в YooKassa
// создает запись о платеже в базе и отправляет запрос в YooKassa
func (s *PaymentService) CreateYooKassaPayment(userID string, courseID string, amount string) (string, error) {
	// проверяем что id курса валидный
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return "", errors.New("invalid course ID")
	}

	// проверяем что id пользователя валидный
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", errors.New("invalid user ID")
	}

	// получаем цену курса из базы
	coursePrice, err := s.paymentRepo.GetCoursePrice(courseUUID)
	if err != nil {
		return "", errors.New("could not get course price")
	}

	// форматируем цену для YooKassa
	priceStr := fmt.Sprintf("%.2f", coursePrice.Amount)

	// создаем запись о платеже в базе
	payment := &entity.Payment{
		PaymentID:    uuid.New(),
		UserID:       userUUID,
		CourseID:     courseUUID,
		Amount:       coursePrice.Amount,
		CurrencyCode: coursePrice.CurrencyCode,
		Status:       "pending",
	}

	// сохраняем платеж в базе
	err = s.paymentRepo.CreatePayment(payment)
	if err != nil {
		return "", errors.New("could not create payment record")
	}

	// формируем запрос для YooKassa
	reqData := &dto.PaymentRequest{}
	reqData.Amount.Value = priceStr
	reqData.Amount.Currency = coursePrice.CurrencyCode
	reqData.Capture = true
	reqData.Description = "Покупка курса"

	reqData.Confirmation.Type = "redirect"
	reqData.Confirmation.ReturnURL = "https://mzt-study.ru/"

	// добавляем метаданные для вебхука
	reqData.Metadata = map[string]string{
		"user_id":    userID,
		"course_id":  courseID,
		"payment_id": payment.PaymentID.String(),
	}

	// генерируем ключ идемпотентности
	idempotenceKey := uuid.New().String()

	// отправляем запрос в YooKassa
	body, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(body))

	req.SetBasicAuth(s.config.Equiring.StoreCode, s.config.Equiring.StoreSecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", idempotenceKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	// проверяем ответ от YooKassa
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(body))
	}

	defer resp.Body.Close()

	// парсим ответ от YooKassa
	var result dto.PaymentResponse
	json.NewDecoder(resp.Body).Decode(&result)

	// обновляем статус платежа в базе
	err = s.paymentRepo.UpdatePaymentStatus(payment.PaymentID, "processing")
	if err != nil {
		fmt.Printf("Error updating payment status: %v\n", err)
	}

	// сохраняем ссылку на платеж в YooKassa
	payment.PaymentRef = result.ID
	err = s.paymentRepo.UpdatePaymentStatus(payment.PaymentID, "processing")
	if err != nil {
		fmt.Printf("Error setting payment reference: %v\n", err)
	}

	// возвращаем ссылку для оплаты
	return result.Confirmation.ConfirmationURL, nil
}

// устанавливает цену для курса
// создает или обновляет запись о цене курса в базе
func (s *PaymentService) SetCoursePrice(courseID uuid.UUID, amount float64, currency string) error {
	// проверяем что курс существует
	_, err := s.courseRepo.GetCourse(courseID)
	if err != nil {
		return errors.New("course not found")
	}

	// проверяем есть ли уже цена для курса
	_, err = s.paymentRepo.GetCoursePrice(courseID)
	if err == nil {
		// если есть, обновляем
		return s.paymentRepo.UpdateCoursePrice(courseID, amount)
	}

	// если нет, создаем новую
	price := &entity.CoursePrice{
		CourseID:     courseID,
		Amount:       amount,
		CurrencyCode: currency,
	}
	return s.paymentRepo.SetCoursePrice(price)
}

// получает цену курса
// просто берет цену курса из базы
func (s *PaymentService) GetCoursePrice(courseID uuid.UUID) (*entity.CoursePrice, error) {
	return s.paymentRepo.GetCoursePrice(courseID)
}

// обновляет статус платежа
// просто меняет статус платежа в базе
func (s *PaymentService) UpdatePaymentStatus(paymentID uuid.UUID, status string) error {
	return s.paymentRepo.UpdatePaymentStatus(paymentID, status)
}

// получает историю платежей пользователя
// берет все платежи пользователя из базы и преобразует их в формат для response
func (s *PaymentService) GetUserTransactions(userID uuid.UUID) ([]dto.PaymentDto, error) {
	// получаем все платежи пользователя
	payments, err := s.paymentRepo.GetPaymentsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// преобразуем каждый платеж в корректный формат для response
	result := make([]dto.PaymentDto, 0, len(payments))
	for _, payment := range payments {
		result = append(result, dto.PaymentDto{
			PaymentID:    payment.PaymentID,
			UserID:       payment.UserID,
			CourseID:     payment.CourseID,
			Amount:       payment.Amount,
			CurrencyCode: payment.CurrencyCode,
			Date:         payment.CreatedAt,
			Status:       payment.Status,
			PaymentRef:   payment.PaymentRef,
		})
	}

	return result, nil
}
