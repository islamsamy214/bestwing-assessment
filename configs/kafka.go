package configs

import "os"

type KafkaCongig map[string]interface{}

func NewKafkaConfig() *KafkaCongig {
	return &KafkaCongig{
		"brokers":      os.Getenv("KAFKA_BROKERS"),
		"group":        os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
		"offset_reset": os.Getenv("KAFKA_OFFSET_RESET"),
	}
}
