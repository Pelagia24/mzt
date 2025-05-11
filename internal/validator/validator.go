package validator

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *Validator) IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func (v *Validator) IsValidTelegram(telegram string) bool {
	const telegramRegex = `^@[A-Za-z0-9_]{5,32}$`
	re := regexp.MustCompile(telegramRegex)
	return re.MatchString(telegram)
}

// TODO adminka
func (v *Validator) IsValidPhoneNumber(phone string) bool {
	const phoneRegex = `^\+?\d{10,15}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

func (v *Validator) IsValidPassword(passwd string) bool {
	if len(passwd) < 8 {
		return false
	}

	allowedChars := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()\[\]\-_=+{}|;:,.<>?/]+$`).MatchString(passwd)
	return allowedChars
}

func (v *Validator) IsValidName(name string) bool {
	validName := regexp.MustCompile(`^[A-Za-zA-Яа-яЁё]{2,}$`)
	return validName.MatchString(name)
}

func (v *Validator) GenerateToken(email, secret string, expirationTimeUnix time.Duration) (string, error) {
	if email == "" {
		return "", errors.New("empty email")
	}
	claims := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(expirationTimeUnix).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (v *Validator) ValidateToken(tokenString, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}
