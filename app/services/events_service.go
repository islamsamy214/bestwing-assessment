package services

import "web-app/app/models/event"

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
