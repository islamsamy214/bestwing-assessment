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

	// Create a channel to handle delivery events
	deliveryChan := make(chan kafka.Event, 10000)

	// Goroutine to handle delivery reports
	go func() {
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
	}()

	// Produce messages to topic
	for i := 0; i < 5; i++ {
		message := fmt.Sprintf(`{"name": "Event %d", "date": "2021-01-01", "user_id": 1}`, i)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, deliveryChan)
	}

	// Flush producer to ensure all messages are delivered
	p.Flush(15 * 1000)

	// Close the delivery channel after flushing
	close(deliveryChan)

	fmt.Println("Produced messages to topic:", topic)
}
