package notification_service

import (
	"errors"
	"fmt"
)

func ConstructSignupEmailVerificationRedisKey(username string, email string, code string) (string, error) {
	if username == "" {
		return "", errors.New("empty username")
	}
	if email == "" {
		return "", errors.New("empty email")
	}
	if code == "" {
		return "", errors.New("empty code")
	}

	return fmt.Sprintf("signup:email-verification:%s:%s:%s", username, email, code), nil
}
