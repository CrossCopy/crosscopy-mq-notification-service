package notification_service

import (
	"errors"
	"github.com/CrossCopy/crosscopy-mq-notification-service/kafka"
	"github.com/CrossCopy/crosscopy-mq-notification-service/redis"
	"time"
)

func SignupHandler(record kafka.SignupTopicRecordValue) error {
	key, err := ConstructSignupEmailVerificationRedisKey(record.Username, record.Email, GetRandomDigitCode(6))
	if err != nil {
		return errors.New("failed to construct signup email verification redis key")
	}
	instance := redis.GetRedisInstance()
	rdb := instance.Client
	ctx := instance.Ctx

	rdb.HSet(ctx, key, "status", "not-verified")
	rdb.HSet(ctx, key, "chance-left", "2")
	rdb.Expire(ctx, key, 10*time.Minute)
	return nil
}
