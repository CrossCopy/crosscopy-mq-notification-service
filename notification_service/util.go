package notification_service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/joho/godotenv"
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
	MailServerHost        string
	KakfaGroupId          string
	KafkaMode             string
	KafkaBootstrapServers string
	KafkaSecurityProtocol string
	KafkaSaslMechanisms   string
	KafkaSaslUsername     string
	KafkaSaslPassword     string
}

func (ev *EnvVars) LoadDotEnv() *EnvVars {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Warning: .env doesn't exist")
	}
	return ev
}

func (ev *EnvVars) LoadEnvVars() *EnvVars {
	host := strings.Split(os.Getenv("MAIL_SERVER"), ":")[0]

	ev.EmailAddress = os.Getenv("EMAIL_ADDRESS")
	ev.EmailPassword = os.Getenv("EMAIL_PASSWORD")
	ev.MailServerAddress = os.Getenv("MAIL_SERVER")
	ev.MailServerHost = host
	ev.KafkaMode = os.Getenv("KAFKA_MODE")
	ev.KakfaGroupId = os.Getenv("KAFKA_GROUP_ID")
	ev.KafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	ev.KafkaSecurityProtocol = os.Getenv("KAFKA_SECURITY_PROTOCOL")
	ev.KafkaSaslMechanisms = os.Getenv("KAFKA_SASL_MECHANISMS")
	ev.KafkaSaslUsername = os.Getenv("KAFKA_SASL_USERNAME")
	ev.KafkaSaslPassword = os.Getenv("KAFKA_SASL_PASSWORD")
	return ev
}

func GetRandomDigitCode(length int) string {
	randomInts := make([]*big.Int, 6)
	for i := 0; i < length; i++ {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(10))
		randomInts[i] = nBig
	}
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(randomInts)), ""), "[]")
}
