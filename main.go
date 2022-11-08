package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CrossCopy/crosscopy-mq-notification-service/notification_service"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RecordValue struct {
	email    string
	username string
}

func main1() {
	err := godotenv.Load()
	if err != nil {
		panic("Error reading .env")
	}
	env := notification_service.LoadEnv()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.6.6.200:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()

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
			recordKey := string(msg.Key)
			recordValue := msg.Value
			record := RecordValue{}
			err = json.Unmarshal(recordValue, &record)
			if err != nil {
				fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
				continue
			}
			key, err := notification_service.ConstructSignupEmailVerificationRedisKey(record.username, record.email, "1024")

			rdb.HSet(ctx, key, "status", "not-verified")
			rdb.HSet(ctx, key, "chance-left", "2")
			rdb.Expire(ctx, key, 10*time.Minute)

			val, err := rdb.Get(ctx, "key").Result()
			fmt.Println(val)
			fmt.Printf("Consumed record with key %s and value %s \n", recordKey, recordValue)
		}
	}

	fmt.Printf("Closing consumer\n")
	consumer.Close()

}
