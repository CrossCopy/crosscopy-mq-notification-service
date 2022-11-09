package notification_service

import (
	"errors"
	"time"

	"github.com/CrossCopy/crosscopy-mq-notification-service/kafka"
	"github.com/CrossCopy/crosscopy-mq-notification-service/redis"
)

func SignupHandler(record kafka.SignupTopicRecordValue) error {
	emailNotifier := GetEmailNotifierInstance()
	verificationCode := GetRandomDigitCode(6)
	key, err := ConstructSignupEmailVerificationRedisKey(record.Username, record.Email, verificationCode)
	if err != nil {
		return errors.New("failed to construct signup email verification redis key")
	}
	instance := redis.GetRedisInstance()
	rdb := instance.Client
	ctx := instance.Ctx
	emailNotifier.SendSignupVerificationEmail(record.Email, verificationCode)
	rdb.HSet(ctx, key, "status", "not-verified")
	rdb.HSet(ctx, key, "chance-left", "2")
	rdb.Expire(ctx, key, 10*time.Minute)
	return nil
}
