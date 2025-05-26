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

		// Parse UUIDs
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
			// Update payment status to succeeded
			err = r.paymentService.UpdatePaymentStatus(paymentIDParsed, "succeeded")
			if err != nil {
				// Log the error but continue with course assignment
			}
		}

		// Assign user to course
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

func (r *Router) CreateCoursePayment(c *gin.Context) {
	courseId := c.Param("course_id")

	user, ok := c.Get("self")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can't get user ID"})
		return
	}

	userId := user.(uuid.UUID)

	// Using empty amount since it will be retrieved from the database
	result, err := r.paymentService.CreateYooKassaPayment(userId.String(), courseId, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment initiated successfully",
		"url":     result,
	})
}
