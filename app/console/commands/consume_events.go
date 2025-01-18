package commands

import (
	"fmt"
	"log"
	"os"
	"web-app/configs"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeEvents() {
	// Kafka configuration
	kafkaConfig := *configs.NewKafkaConfig()
	topic := os.Getenv("KAFKA_EVENTS_TOPIC")

	// Create a Kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaConfig["brokers"].(string),
		"group.id":           kafkaConfig["group"].(string),
		"auto.offset.reset":  os.Getenv("KAFKA_OFFSET_RESET"), // Consume messages from the beginning if no offset is stored
		"enable.auto.commit": false,                           // Disable auto commit for manual offset management
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	// Subscribe to the topic
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	fmt.Printf("Consumer started, waiting for messages from topic %s...\n", topic)

	// Consume messages in a loop
	for {
		msg, err := consumer.ReadMessage(-1) // -1: Block indefinitely until a message is received
		if err != nil {
			// Handle error, e.g., timeout or other issues
			if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrTimedOut {
				continue
			}
			log.Printf("Error consuming message: %s\n", err)
			break
		}

		// Process the message
		fmt.Printf("Received message: Key=%s, Value=%s, Topic=%s, Partition=%d, Offset=%d\n",
			string(msg.Key), string(msg.Value), *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)

		// Commit the offset manually
		_, err = consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Error committing offset: %s\n", err)
		}
	}
}
