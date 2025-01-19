package commands

import (
	"fmt"
	"log"
	"os"
	"web-app/configs"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ProduceEvents() {
	// Kafka configuration
	topic := os.Getenv("KAFKA_EVENTS_TOPIC")

	// Create a new producer
	p, err := kafka.NewProducer(configs.NewKafkaProducerConfig())
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}

	defer p.Close()

	// Produce messages to topic (asynchronously)
	deliveryChan := make(chan kafka.Event, 10000)
	for i := 0; i < 30; i++ {
		message := fmt.Sprintf(`{"name": "Event %d", "date": "2021-01-01", "user_id": 1}`, i)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, deliveryChan)
	}

	// Wait for message deliveries
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
			}
		}
	}

	close(deliveryChan)

	// Flush the producer before closing
	p.Flush(15 * 1000)

	fmt.Println("Produced 100 messages to topic", topic)
}
