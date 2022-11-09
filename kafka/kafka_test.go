package kafka

import (
	"fmt"
	"testing"
)

var TestKafkaConfig = KafkaConfig{
	Mode:             LocalMode,
	BootstrapServers: "localhost:9093",
	SecurityProtocol: "SASL_PLAINTEXT",
	SaslMechanisms:   "PLAIN",
	SaslUsername:     "",
	SaslPassword:     "",
	GroupId:          "crosscopy-notification",
}

func TestCreateTopic(t *testing.T) {
	adminClient := GetAdminClient(TestKafkaConfig)
	topic := "some_topic5"
	CreateTopic(adminClient, topic, 1)
	_, topics := GetMetadata(adminClient, &topic)
	for _, topic_ := range topics {
		if topic_ == topic {
			return
		}
	}
	t.Errorf("Topic %s not deleted\n", topic)
	DeleteTopic(adminClient, []string{topic})
	adminClient.Close()
}

func TestDeleteTopic(t *testing.T) {
	topic := "some_topic"
	adminClient := GetAdminClient(TestKafkaConfig)
	replicationFactor := GetReplicationFactor(TestKafkaConfig.Mode)
	CreateTopic(adminClient, topic, replicationFactor)
	DeleteTopic(adminClient, []string{topic})
	_, topics := GetMetadata(adminClient, &topic)
	for _, topic_ := range topics {
		if topic_ == topic {
			t.Errorf("Topic %s not deleted\n", topic)
		}
	}
	adminClient.Close()
}

func TestProduceSignupTopicMessage(t *testing.T) {
	producer := GetProducer(TestKafkaConfig)
	adminClient := GetAdminClientFromProducer(producer)
	ListenProducerEvents(producer)
	topic := "signup"
	CreateTopic(adminClient, topic, 1)
	for i := 10; i < 20; i++ {
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@crosscopy.io", i)
		err := CreateSignupTopicMessage(producer, username, email)
		if err != nil {
			fmt.Printf("Failed to create consumer: %s", err)
		}
	}
	// _, topics := GetMetadata(producer, &topic)
	// fmt.Println(topics)
	producer.Close()
}
