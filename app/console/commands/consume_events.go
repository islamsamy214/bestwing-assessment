package commands

import (
	"fmt"
	"log"
	"os"
	"web-app/app/services"
	"web-app/configs"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var isConsumerRunning = false

func ConsumeEvents() {
	// Kafka configuration
	topic := os.Getenv("KAFKA_EVENTS_TOPIC")

	// Create a Kafka consumer
	consumer, err := kafka.NewConsumer(configs.NewKafkaConsumerConfig())
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	// Subscribe to the topic
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	// Start the consumer
	if isConsumerRunning {
		log.Println("Kafka consumer already running.")
		return
	}
	isConsumerRunning = true
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
		log.Printf("Received message: Key=%s, Value=%s, Topic=%s, Partition=%d, Offset=%d\n",
			string(msg.Key), string(msg.Value), *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)

		if err := services.NewEventService().HandleEventConsumption(msg.Value); err != nil {
			log.Printf("Error handling event consumption: %s\n", err)
			continue
		}

		// Commit the offset manually
		_, err = consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Error committing offset: %s\n", err)
		}
	}
}
