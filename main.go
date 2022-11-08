package main

import (
	"fmt"
	"github.com/confluentinc/examples/clients/cloud/go/ccloud"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
)

func main() {
	configFile, topic := ccloud.ParseArgs()
	conf := ccloud.ReadCCloudConfig(*configFile)
	// Create Consumer instance
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf["bootstrap.servers"],
		"sasl.mechanisms":   conf["sasl.mechanisms"],
		"security.protocol": conf["security.protocol"],
		"sasl.username":     conf["sasl.username"],
		"sasl.password":     conf["sasl.password"],
		"group.id":          "go_example_group_1",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	// Subscribe to topic
	err = consumer.SubscribeTopics([]string{*topic}, nil)

	//ret, err := notification_service.ConstructSignupEmailVerificationRedisKey("dev", "dev@crosscopy.io", "1024")
	//ret, err = notification_service.ConstructSignupEmailVerificationRedisKey("", "dev@crosscopy.io", "1024")
	//if err != nil {
	//	panic(err)
	//}
	//println(ret)

}
