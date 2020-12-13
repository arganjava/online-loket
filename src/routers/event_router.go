package routers

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/service_repository"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
	"net/http"
)

type EventRouter struct {
	eventRepository service_repository.EventRepository
}

func NewEventRouter(eventRepository service_repository.EventRepository) *EventRouter {
	return &EventRouter{
		eventRepository: eventRepository,
	}
}

func (l EventRouter) CreateEvent(c *gin.Context) {
	var request dto.EventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := validator.Validate(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	_, err := l.eventRepository.CreateEvent(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Event created successfully!"})

}
