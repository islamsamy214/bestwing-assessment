package services

import (
	"encoding/json"
	"web-app/app/models/event"
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

	return nil
}
