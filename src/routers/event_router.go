package routers

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/interfaces"
	"github.com/arganjava/online-loket/src/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
	"net/http"
)

type EventRouter struct {
	eventRepository interfaces.EventRepository
}

func NewEventRouter(eventRepository interfaces.EventRepository) *EventRouter {
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

func (l EventRouter) GetEventInfo(c *gin.Context) {
	id := c.Param("id")
	data, err := l.eventRepository.FindEventById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	} else if data == nil {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Data not found", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Event Get successfully!", "data": buildEventResponse(data)})

}

func (l EventRouter) CreateEventTicket(c *gin.Context) {
	var request dto.EventTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if request.Type != "ADULT" && request.Type != "CHILD" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Type must between CHILD or ADULT"})
		return
	}

	_, err := l.eventRepository.CreateEventTicket(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Event Ticket created successfully!"})
}

func buildEventResponse(data *models.Event) dto.EventResponse {

	return dto.EventResponse{
		ID:            data.ID,
		EventName:     data.Name,
		Description:   data.Description,
		ScheduleBegin: data.ScheduleBegin.String(),
		ScheduleEnd:   data.ScheduleEnd.String(),
		Location:      buildLocationResponse(data.Location),
		EventTickets:  buildEventTicketsResponse(data.EventTickets),
	}
}

func buildEventTicketsResponse(tickets []*models.EventTicket) []dto.EventTicketResponse {
	eventTicketResponses := make([]dto.EventTicketResponse, 0)
	for _, data := range tickets {
		eventTicketResponses = append(eventTicketResponses, buildEventTicketResponse(data))
	}
	return eventTicketResponses
}

func buildEventTicketResponse(data *models.EventTicket) dto.EventTicketResponse {
	return dto.EventTicketResponse{
		ID:       data.ID,
		Type:     data.Type,
		Price:    data.Price,
		Quantity: data.Quantity,
	}
}

func buildLocationResponse(location *models.Location) dto.LocationResponse {
	return dto.LocationResponse{
		ID:       location.ID,
		Village:  location.Village,
		CityName: location.CityName,
		Address:  location.Address,
		Country:  location.Country,
	}
}
