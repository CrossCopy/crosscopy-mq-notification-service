package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const (
	TopicSignup string = "signup"
)

type SignupTopicRecordValue struct {
	Email    string
	Username string
}

type KafkaMode string

const (
	LocalMode KafkaMode = "local"
	CloudMode KafkaMode = "cloud"
)

type KafkaConfig struct {
	Mode             KafkaMode
	BootstrapServers string
	SaslMechanisms   string
	SecurityProtocol string
	SaslUsername     string
	SaslPassword     string
	GroupId          string
}

func GetReplicationFactor(mode KafkaMode) int {
	replicationFactor := 1
	if mode == CloudMode {
		replicationFactor = 3
	}
	return replicationFactor
}

func GetAdminClientFromProducer(producer *kafka.Producer) *kafka.AdminClient {
	adminClient, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		fmt.Printf("Failed to create new admin client from producer: %s", err)
		os.Exit(1)
	}
	return adminClient
}

func GetAdminClient(config KafkaConfig) *kafka.AdminClient {
	configMap := GetKafkaConfigMap(config)
	adminClient, err := kafka.NewAdminClient(configMap)
	if err != nil {
		fmt.Printf("Failed to create new admin client from producer: %s", err)
		os.Exit(1)
	}
	return adminClient
}

// CreateTopic creates a topic using the Admin Client API
// CreateTopic(producer, topic)
func CreateTopic(adminClient *kafka.AdminClient, topic string, replicationFactor int) {
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Printf("ParseDuration(60s): %s", err)
		os.Exit(1)
	}

	results, err := adminClient.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: replicationFactor}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		fmt.Printf("Admin Client request error: %v\n", err)
		os.Exit(1)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			fmt.Printf("Failed to create topic: %v\n", result.Error)
			os.Exit(1)
		}
		fmt.Printf("%v\n", result)
	}
}

func DeleteTopic(adminClient *kafka.AdminClient, topics []string) {
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Printf("ParseDuration(60s): %s", err)
		os.Exit(1)
	}

	results, err := adminClient.DeleteTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		topics,
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		fmt.Printf("Admin Client request error: %v\n", err)
		os.Exit(1)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			fmt.Printf("Failed to create topic: %v\n", result.Error)
			os.Exit(1)
		}
		fmt.Printf("%v\n", result)
	}
}

func GetKafkaConfigMap(config KafkaConfig) *kafka.ConfigMap {
	cm := &kafka.ConfigMap{"bootstrap.servers": config.BootstrapServers}
	if config.Mode == CloudMode {
		cm.SetKey("sasl.mechanisms", config.SaslMechanisms)
		cm.SetKey("security.protocol", config.SecurityProtocol)
		cm.SetKey("sasl.username", config.SaslUsername)
		cm.SetKey("sasl.password", config.SaslPassword)
	} else if config.Mode == LocalMode {
		// pass
	} else {
		fmt.Printf("Wrong kafka config mode: %s\n", config.Mode)
	}
	return cm
}

func GetProducer(config KafkaConfig) *kafka.Producer {
	configMap := GetKafkaConfigMap(config)
	fmt.Println("configMap")
	fmt.Println(configMap)

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}
	return producer
}

func GetConsumer(config KafkaConfig) *kafka.Consumer {
	configMap := GetKafkaConfigMap(config)
	configMap.SetKey("group.id", config.GroupId)
	configMap.SetKey("auto.offset.reset", "earliest")
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	return consumer
}

func GetMetadata(adminClient *kafka.AdminClient, topic *string) (*kafka.Metadata, []string) {
	metadata, err := adminClient.GetMetadata(topic, true, 100000)
	if err != nil {
		fmt.Printf("Failed to get metadata: %s", err)
		os.Exit(1)
	}
	topics := make([]string, len(metadata.Topics))

	i := 0
	for k := range metadata.Topics {
		topics[i] = k
		i++
	}
	return metadata, topics
}

func ListenProducerEvents(producer *kafka.Producer) {
	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()
}

func CreateSignupTopicMessage(producer *kafka.Producer, username string, email string) error {

	topic := "signup"

	//for n := 0; n < 10; n++ {
	recordKey := "signup"
	data := &SignupTopicRecordValue{
		Username: username,
		Email:    email,
	}
	recordValue, _ := json.Marshal(&data)
	fmt.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(recordKey),
		Value:          []byte(recordValue),
	}, nil)
	//}

	// Wait for all messages to be delivered
	producer.Flush(15 * 1000)
	return err
}
