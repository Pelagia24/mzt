package router

// // filepath: c:\Users\mamev\mzt-all\mzt\internal\router\router_test.go
// package router

// import (
// 	"mzt/config"
// 	"mzt/internal/middleware"
// 	"mzt/internal/service"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// func TestNewRouter(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	handler := gin.New()

// 	authService := &service.UserService{}
// 	courseService := &service.CourseService{}
// 	paymentService := &service.PaymentService{}
// 	mw := &middleware.Middleware{}

// 	cfg := &config.Config{}
// 	router := NewRouter(cfg, handler, authService, courseService, mw)

// 	assert.NotNil(t, router)
// 	assert.NotNil(t, router.authService)
// 	assert.NotNil(t, router.courseService)
// 	assert.NotNil(t, router.config)
// }

// func TestPaymentWebhookRoute(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	handler := gin.New()

// 	authService := &service.UserService{}
// 	courseService := &service.CourseService{}
// 	paymentService := &service.PaymentService{}
// 	mw := &middleware.Middleware{}

// 	cfg := &config.Config{
// 		Equiring: config.Equiring{
// 			SecretPath: "/test-secret",
// 		},
// 	}
// 	NewRouter(cfg, handler, authService, courseService, mw)

// 	w := performRequest(handler, "POST", "/api/v1/webhook/payments/test-secret", nil)
// 	assert.Equal(t, 404, w.Code) // Assuming no handler logic is implemented yet
// }