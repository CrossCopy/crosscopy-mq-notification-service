package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mykafka "github.com/CrossCopy/crosscopy-mq-notification-service/kafka"
	"github.com/CrossCopy/crosscopy-mq-notification-service/notification_service"
	myredis "github.com/CrossCopy/crosscopy-mq-notification-service/redis"
)

func main() {
	env := notification_service.EnvVars{}
	env.LoadDotEnv().LoadEnvVars()
	notification_service.GetEmailNotifierInstance().Init(env.EmailAddress, env.EmailPassword, env.MailServerAddress, env.MailServerHost)
	myredis.GetRedisInstance().Connect(env.REDIS_HOST, env.REDIS_PASS, env.REDIS_PORT)
	redisInstance := myredis.GetRedisInstance()
	rdb := redisInstance.Client
	ctx := redisInstance.Ctx
	if err := rdb.Del(ctx, "test-connection").Err(); err != nil {
		panic(err)
	}
	if err := rdb.Set(ctx, "test-connection", "huakun", 10*time.Minute).Err(); err != nil {
		panic(err)
	}
	textConnectionVal, err := rdb.Get(ctx, "test-connection").Result()
	if err != nil {
		panic(err)
	}
	if textConnectionVal != "huakun" {
		panic("Redis Connection Not Successful")
	}
	// Create Consumer instance
	var kafkaConfig = mykafka.KafkaConfig{
		Mode:             mykafka.KafkaMode(env.KafkaMode),
		BootstrapServers: env.KafkaBootstrapServers,
		SecurityProtocol: env.KafkaSecurityProtocol,
		SaslMechanisms:   env.KafkaSaslMechanisms,
		SaslUsername:     env.KafkaSaslUsername,
		SaslPassword:     env.KafkaSaslPassword,
		GroupId:          env.KakfaGroupId,
	}
	adminClient := mykafka.GetAdminClient(kafkaConfig)
	replicationFactor := mykafka.GetReplicationFactor(mykafka.KafkaMode(env.KafkaMode))
	mykafka.CreateTopic(adminClient, mykafka.TopicSignup, replicationFactor)
	consumer := mykafka.GetConsumer(kafkaConfig)
	topicsToSubscribe := []string{mykafka.TopicSignup}
	// Subscribe to topic
	fmt.Println("Consumer Starts")
	err1 := consumer.SubscribeTopics(topicsToSubscribe, nil)
	if err1 != nil {
		panic("Error Subscribing to Topic")
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			//recordKey := string(msg.Key)
			switch *msg.TopicPartition.Topic {

			case mykafka.TopicSignup:
				fmt.Printf("Topic Received: %s\n", mykafka.TopicSignup)
				recordValue := msg.Value
				record := mykafka.SignupTopicRecordValue{}
				err = json.Unmarshal(recordValue, &record)
				if err != nil {
					fmt.Printf("Failed to decode JSON at offset %d: %v\n", msg.TopicPartition.Offset, err)
					continue
				}
				if err := notification_service.SignupHandler(record); err != nil {
					fmt.Println(err)
				}
			default:
				fmt.Printf("error: unhandled topic %s\n", *msg.TopicPartition.Topic)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	if err := consumer.Close(); err != nil {
		panic(err)
	}

}
