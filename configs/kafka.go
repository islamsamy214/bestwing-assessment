package configs

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewKafkaConsumerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":     os.Getenv("KAFKA_BROKERS"),
		"group.id":              os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
		"auto.offset.reset":     os.Getenv("KAFKA_OFFSET_RESET"),
		"enable.auto.commit":    false,
		"session.timeout.ms":    30000,
		"max.poll.interval.ms":  300000,
		"heartbeat.interval.ms": 3000,
		"socket.timeout.ms":     30000,
	}
}

func NewKafkaProducerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
		"socket.timeout.ms": 60000,
	}
}
