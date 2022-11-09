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

func TestDummy(t *testing.T) {
	fmt.Println(TestKafkaConfig)
	fmt.Println(TestKafkaConfig.SaslUsername)
	fmt.Println(TestKafkaConfig.SaslUsername == "")
}
func TestCreateTopic(t *testing.T) {
	producer := GetProducer(TestKafkaConfig)
	topic := "some_topic1"
	CreateTopic(producer, topic, 1)
	_, topics := GetMetadata(producer, &topic)
	for _, topic_ := range topics {
		if topic_ == topic {
			return
		}
	}
	t.Errorf("Topic %s not deleted\n", topic)
	DeleteTopic(producer, []string{topic})
	producer.Close()
}

func TestDeleteTopic(t *testing.T) {
	producer := GetProducer(TestKafkaConfig)
	topic := "some_topic"
	CreateTopic(producer, topic, 1)
	DeleteTopic(producer, []string{topic})
	_, topics := GetMetadata(producer, &topic)
	for _, topic_ := range topics {
		if topic_ == topic {
			t.Errorf("Topic %s not deleted\n", topic)
		}
	}
	producer.Close()
}

func TestProduceSignupTopicMessage(t *testing.T) {
	producer := GetProducer(TestKafkaConfig)
	ListenProducerEvents(producer)
	topic := "signup"
	CreateTopic(producer, topic, 1)
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
