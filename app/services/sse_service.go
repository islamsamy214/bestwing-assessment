package services

import "fmt"

type SSEService struct {
	Channel chan map[string]interface{}
}

var instance *SSEService

func NewSSEService() *SSEService {
	if instance == nil {
		instance = &SSEService{
			Channel: make(chan map[string]interface{}, 100), // Buffered channel
		}
	}
	return instance
}

// SendEvent adds an event to the channel and logs it
func (s *SSEService) SendEvent(event map[string]interface{}) {
	if s == nil {
		fmt.Println("Error: SSEService is nil")
		return
	}
	fmt.Printf("Event added to channel: %+v\n", event)
	s.Channel <- event
	fmt.Printf("Server Channel address: %p\n", s.Channel)
}
