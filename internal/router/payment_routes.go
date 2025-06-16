package router

import (
	"net/http"

	"mzt/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Router) YooWebhookHandler(c *gin.Context) {
	var webhook dto.YooWebhook

	if err := c.BindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook"})
		return
	}

	if webhook.Event == "payment.succeeded" && webhook.Object.Status == "succeeded" {
		userID := webhook.Object.Metadata.UserID
		courseID := webhook.Object.Metadata.CourseID
		paymentID := webhook.Object.Metadata.PaymentID

		// преобразуем строки в uuid
		userIDParsed, err := uuid.Parse(userID)
		if err != nil {
			//TODO handle this(i think log + notification somewhere)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		courseIDParsed, err := uuid.Parse(courseID)
		if err != nil {
			//TODO handle this(i think log + notification somewhere)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		paymentIDParsed, err := uuid.Parse(paymentID)
		if err != nil {
			// Log but continue without payment update
			// We should still assign the course
		} else {
			// обновляем статус платежа на успешный
			err = r.paymentService.UpdatePaymentStatus(paymentIDParsed, "succeeded")
			if err != nil {
				// Log the error but continue with course assignment
			}
		}

		// записываем пользователя на курс
		err = r.courseService.AssignUserToCourse(courseIDParsed, userIDParsed)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
}

// создает платеж для курса
// возвращает ссылку на оплату
func (r *Router) CreateCoursePayment(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")

	// достаем id пользователя из контекста
	user, ok := c.Get("self")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can't get user ID"})
		return
	}

	userId := user.(uuid.UUID)

	// создаем платеж через сервис
	// сумма будет получена из базы данных
	result, err := r.paymentService.CreateYooKassaPayment(userId.String(), courseId, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// отправляем ссылку на оплату клиенту
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment initiated successfully",
		"url":     result,
	})
}

// получает список транзакций пользователя
// доступно только админам
func (r *Router) GetUserTransactions(c *gin.Context) {
	// достаем id пользователя из параметров запроса
	userID := c.Param("user_id")

	// преобразуем строку в uuid
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// получаем список транзакций из сервиса
	transactions, err := r.paymentService.GetUserTransactions(userIDParsed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// отправляем список транзакций клиенту
	c.JSON(http.StatusOK, transactions)
}
