package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"web-app/app/console/commands"
	"web-app/app/models/event"
	"web-app/app/services"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventsService *services.EventService
	sseService    *services.SSEService
}

func NewEventController() *EventController {
	return &EventController{
		eventsService: services.NewEventService(),
		sseService:    services.NewSSEService(),
	}
}

func (e *EventController) Index(c *gin.Context) {

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "1"))

	events, err := e.eventsService.Index(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": events})
}

func (e *EventController) Create(c *gin.Context) {
	eventsModel := event.NewEventModel()
	if err := c.ShouldBindJSON(eventsModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	eventsModel.UserId = c.MustGet("userId").(int64)
	if err := services.NewEventService().ProduceEvent(eventsModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": eventsModel})
}

func (e *EventController) Listen(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		http.Error(c.Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Start the kafka consumer
	go commands.ConsumeEvents()

	log.Printf("Listener Channel address: %p\n", e.sseService.Channel)

	// Listen for events
	for {
		select {
		case event := <-e.sseService.Channel:
			jsonData, err := json.Marshal(event)
			if err != nil {
				// Log the error for debugging
				c.Writer.Write([]byte("data: {\"error\": \"Failed to serialize event\"}\n\n"))
			} else {
				log.Printf("Sending event: %s", string(jsonData))
				c.Writer.Write([]byte("data: " + string(jsonData) + "\n\n"))
			}
			flusher.Flush()
		case <-c.Writer.CloseNotify():
			return // Exit on client disconnect
		}
	}
}
