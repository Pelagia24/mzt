package dto

type PaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Capture     bool              `json:"capture"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

type PaymentResponse struct {
	ID           string `json:"id"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

type YooWebhook struct {
	Event  string `json:"event"`
	Object struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Metadata struct {
			UserID    string `json:"user_id"`
			CourseID  string `json:"course_id"`
			PaymentID string `json:"payment_id"`
		} `json:"metadata"`
	} `json:"object"`
}
