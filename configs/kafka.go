package configs

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewKafkaConsumerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_BROKERS"),
		"group.id":           os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
		"auto.offset.reset":  os.Getenv("KAFKA_OFFSET_RESET"),
		"enable.auto.commit": false, // Disable auto commit for manual offset management
	}
}

func NewProducerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
	}
}
