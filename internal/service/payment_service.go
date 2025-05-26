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

type PaymentService struct {
	config      *config.Config
	courseRepo  repository.CourseRepository
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(cfg *config.Config, courseRepo repository.CourseRepository, paymentRepo repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		config:      cfg,
		courseRepo:  courseRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *PaymentService) CreateYooKassaPayment(userID string, courseID string, amount string) (string, error) {
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return "", errors.New("invalid course ID")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", errors.New("invalid user ID")
	}

	coursePrice, err := s.paymentRepo.GetCoursePrice(courseUUID)
	if err != nil {
		return "", errors.New("could not get course price")
	}

	priceStr := fmt.Sprintf("%.2f", coursePrice.Amount)

	payment := &entity.Payment{
		PaymentID:    uuid.New(),
		UserID:       userUUID,
		CourseID:     courseUUID,
		Amount:       coursePrice.Amount,
		CurrencyCode: coursePrice.CurrencyCode,
		Status:       "pending",
	}

	err = s.paymentRepo.CreatePayment(payment)
	if err != nil {
		return "", errors.New("could not create payment record")
	}

	reqData := &dto.PaymentRequest{}
	reqData.Amount.Value = priceStr
	reqData.Amount.Currency = coursePrice.CurrencyCode
	reqData.Capture = true
	reqData.Description = "Покупка курса"

	reqData.Confirmation.Type = "redirect"
	reqData.Confirmation.ReturnURL = "https://example.com/example_url"

	reqData.Metadata = map[string]string{
		"user_id":    userID,
		"course_id":  courseID,
		"payment_id": payment.PaymentID.String(),
	}

	idempotenceKey := uuid.New().String()

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

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(body))
	}

	defer resp.Body.Close()

	var result dto.PaymentResponse
	json.NewDecoder(resp.Body).Decode(&result)

	err = s.paymentRepo.UpdatePaymentStatus(payment.PaymentID, "processing")
	if err != nil {
		fmt.Printf("Error updating payment status: %v\n", err)
	}

	payment.PaymentRef = result.ID
	err = s.paymentRepo.UpdatePaymentStatus(payment.PaymentID, "processing")
	if err != nil {
		fmt.Printf("Error setting payment reference: %v\n", err)
	}

	return result.Confirmation.ConfirmationURL, nil
}

func (s *PaymentService) SetCoursePrice(courseID uuid.UUID, amount float64, currency string) error {
	_, err := s.courseRepo.GetCourse(courseID)
	if err != nil {
		return errors.New("course not found")
	}

	_, err = s.paymentRepo.GetCoursePrice(courseID)
	if err == nil {
		return s.paymentRepo.UpdateCoursePrice(courseID, amount)
	}

	price := &entity.CoursePrice{
		CourseID:     courseID,
		Amount:       amount,
		CurrencyCode: currency,
	}
	return s.paymentRepo.SetCoursePrice(price)
}

func (s *PaymentService) GetCoursePrice(courseID uuid.UUID) (*entity.CoursePrice, error) {
	return s.paymentRepo.GetCoursePrice(courseID)
}

func (s *PaymentService) UpdatePaymentStatus(paymentID uuid.UUID, status string) error {
	return s.paymentRepo.UpdatePaymentStatus(paymentID, status)
}
