package notification_service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
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

type EnvVars struct {
	EmailAddress          string
	EmailPassword         string
	MailServerAddress     string
	KakfaGroupId          string
	KafkaMode             string
	KafkaBootstrapServers string
	KafkaSecurityProtocol string
	KafkaSaslMechanisms   string
	KafkaSaslUsername     string
	KafkaSaslPassword     string
}

func LoadEnv() EnvVars {
	env := EnvVars{
		EmailAddress:          os.Getenv("EMAIL_ADDRESS"),
		EmailPassword:         os.Getenv("EMAIL_PASSWORD"),
		MailServerAddress:     os.Getenv("MAIL_SERVER"),
		KafkaMode:             os.Getenv("KAFKA_MODE"),
		KakfaGroupId:          os.Getenv("KAFKA_GROUP_ID"),
		KafkaBootstrapServers: os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		KafkaSecurityProtocol: os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		KafkaSaslMechanisms:   os.Getenv("KAFKA_SASL_MECHANISMS"),
		KafkaSaslUsername:     os.Getenv("KAFKA_SASL_USERNAME"),
		KafkaSaslPassword:     os.Getenv("KAFKA_SASL_PASSWORD"),
	}
	return env
}

func GetRandomDigitCode(length int) string {
	randomInts := make([]*big.Int, 6)
	for i := 0; i < length; i++ {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(10))
		randomInts[i] = nBig
	}
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(randomInts)), ""), "[]")
}
