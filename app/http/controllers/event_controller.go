package controllers

import (
	"net/http"
	"strconv"
	"web-app/app/models/event"
	"web-app/app/services"

	"github.com/gin-gonic/gin"
)

type EventController struct{}

var eventsService *services.EventService

func NewEventController() *EventController {
	eventsService = services.NewEventService()
	return &EventController{}
}

func (e *EventController) Index(c *gin.Context) {

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "1"))

	events, err := eventsService.Index(limit, offset)
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
