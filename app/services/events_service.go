package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"web-app/app/models/event"
	"web-app/configs"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventService struct{}

func NewEventService() *EventService {
	return &EventService{}
}

func (e *EventService) Index(limit, offset int) ([]event.Event, error) {
	eventsModel := event.NewEventModel()
	events, err := eventsModel.Paginate(limit, offset)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (e *EventService) Create(eventsModel *event.Event) error {
	if err := eventsModel.Create(); err != nil {
		return err
	}
	return nil
}

func (e *EventService) ProduceEvent(eventsModel *event.Event) error {
	// Kafka configuration
	topic := os.Getenv("KAFKA_EVENTS_TOPIC")

	// Create a new producer
	producer, err := kafka.NewProducer(configs.NewKafkaProducerConfig())
	if err != nil {
		log.Printf("failed to create producer: %v\n", err)
		return err
	}
	defer producer.Close()

	// Generate a message
	message := fmt.Sprintf(`{"name": "%s", "date": "%s", "user_id": %d}`, eventsModel.Name, eventsModel.Date, eventsModel.UserId)

	// Error channel for delivery report
	errChan := make(chan error, 1)

	// Produce the message to the topic
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Handle the delivery report
	go func() {
		defer close(errChan)
		for ev := range producer.Events() {
			switch msg := ev.(type) {
			case *kafka.Message:
				if msg.TopicPartition.Error != nil {
					errChan <- msg.TopicPartition.Error
				} else {
					errChan <- nil
				}
				return
			}
		}
	}()

	// Wait for the delivery result
	if deliveryErr := <-errChan; deliveryErr != nil {
		log.Printf("failed to deliver message: %v\n", deliveryErr)
		return deliveryErr
	}

	// Flush remaining messages
	producer.Flush(15 * 1000)

	return nil
}

func (e *EventService) HandleEventConsumption(eventData []byte) error {
	// Unmarshal the event data
	var eventDataJson struct {
		Name   string `json:"name"`
		Date   string `json:"date"`
		UserId int64  `json:"user_id"`
	}

	if err := json.Unmarshal(eventData, &eventDataJson); err != nil {
		return err
	}

	// Create a new event
	eventsModel := event.NewEventModel()
	eventsModel.Name = eventDataJson.Name
	eventsModel.Date = eventDataJson.Date
	eventsModel.UserId = eventDataJson.UserId

	if err := e.Create(eventsModel); err != nil {
		return err
	}

	eventPayload := map[string]interface{}{
		"id":      eventsModel.ID,
		"name":    eventsModel.Name,
		"date":    eventsModel.Date,
		"user_id": eventsModel.UserId,
	}

	// Send the event to the websocket server
	sseSevice := NewSSEService()
	sseSevice.Channel <- eventPayload
	sseSevice.SendEvent(eventPayload)

	return nil
}
