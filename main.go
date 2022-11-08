package main

import "github.com/CrossCopy/crosscopy-mq-notification-service/notification_service"

func main() {
	ret, err := notification_service.ConstructSignupEmailVerificationRedisKey("dev", "dev@crosscopy.io", "1024")
	ret, err = notification_service.ConstructSignupEmailVerificationRedisKey("", "dev@crosscopy.io", "1024")
	if err != nil {
		panic(err)
	}
	println(ret)

}
