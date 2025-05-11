package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/repository"
	"net/http"

	"github.com/google/uuid"
)

type PaymentService struct {
	config *config.Config
	repo   repository.CourseRepository
}

func NewPaymentService(cfg *config.Config, repo repository.CourseRepository) *PaymentService {
	return &PaymentService{
		config: cfg,
		repo:   repo,
	}
}

func (s *PaymentService) CreateYooKassaPayment(userID string, courseID string, amount string) (string, error) {
	reqData := &dto.PaymentRequest{}
	//TODO get amount from db
	reqData.Amount.Value = amount
	reqData.Amount.Currency = "RUB"
	reqData.Capture = true
	reqData.Description = "Покупка курса"

	reqData.Confirmation.Type = "redirect"
	reqData.Confirmation.ReturnURL = "https://example.com/example_url"

	reqData.Metadata = map[string]string{
		"user_id":   userID,
		"course_id": courseID,
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
	fmt.Println(resp.Body)
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Confirmation.ConfirmationURL, nil
}
