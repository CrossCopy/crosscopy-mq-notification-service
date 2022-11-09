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
	"github.com/CrossCopy/crosscopy-mq-notification-service/redis"
)

func main() {
	env := notification_service.EnvVars{}
	env.LoadDotEnv().LoadEnvVars()
	notification_service.GetEmailNotifierInstance().Init(env.EmailAddress, env.EmailPassword, env.MailServerAddress, env.MailServerHost)
	redis.GetRedisInstance().Connect()
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
	err := consumer.SubscribeTopics(topicsToSubscribe, nil)
	if err != nil {
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
					fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
					continue
				}
				notification_service.SignupHandler(record)
			default:
				fmt.Printf("error: unhandled topic %s\n", *msg.TopicPartition.Topic)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	consumer.Close()

}
