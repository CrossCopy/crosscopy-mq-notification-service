package notification_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/CrossCopy/crosscopy-mq-notification-service/kafka"
	"github.com/CrossCopy/crosscopy-mq-notification-service/redis"
)

func SignupHandler(record kafka.SignupTopicRecordValue) error {
	env := EnvVars{}
	env.LoadDotEnv().LoadEnvVars()
	emailNotifier := GetEmailNotifierInstance()
	verificationCode := GetRandomDigitCode(6)
	key, err := ConstructSignupEmailVerificationRedisKey(record.Username, record.Email, verificationCode)
	if err != nil {
		return errors.New("failed to construct signup email verification redis key")
	}
	redis.GetRedisInstance().Connect(env.REDIS_HOST, env.REDIS_PASS, env.REDIS_PORT)
	instance := redis.GetRedisInstance()
	rdb := instance.Client
	ctx := instance.Ctx
	if err := rdb.HSet(ctx, key, "status", "not-verified").Err(); err != nil {
		fmt.Println(err)
	}
	if err := rdb.HSet(ctx, key, "chance-left", "2").Err(); err != nil {
		fmt.Println(err)
	}
	if err := rdb.Expire(ctx, key, 10*time.Minute).Err(); err != nil {
		fmt.Println(err)
	}
	emailNotifier.SendSignupVerificationEmail(record.Email, verificationCode)

	return nil
}
