package main

import (
	"encoding/json"
	"fmt"
	main2 "github.com/CrossCopy/crosscopy-mq-notification-service/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CrossCopy/crosscopy-mq-notification-service/notification_service"
	redis2 "github.com/CrossCopy/crosscopy-mq-notification-service/redis"
	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error reading .env")
	}
	env := notification_service.LoadEnv()

	instance := redis2.GetRedisInstance()
	instance.Connect()

	// Create Consumer instance
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": env.KafkaBootstrapServers,
		"sasl.mechanisms":   env.KafkaSaslMechanisms,
		"security.protocol": env.KafkaSecurityProtocol,
		"sasl.username":     env.KafkaSaslUsername,
		"sasl.password":     env.KafkaSaslPassword,
		"group.id":          env.KakfaGroupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	topicsToSubscribe := []string{"test1"}
	// Subscribe to topic
	err = consumer.SubscribeTopics(topicsToSubscribe, nil)
	if err != nil {
		panic("Error Subscribing to Topic")
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run == true {
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
			case "test1":
				fmt.Println("test1")

			case "signup":
				fmt.Println("test1")
				recordValue := msg.Value
				record := main2.SignupTopicRecordValue{}
				err = json.Unmarshal(recordValue, &record)
				if err != nil {
					fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
					continue
				}
			default:
				fmt.Printf("error: unhandled topic %s\n", *msg.TopicPartition.Topic)
			}

		}
	}

	fmt.Printf("Closing consumer\n")
	consumer.Close()

}
