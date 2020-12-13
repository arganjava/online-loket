package dto

type EventRequest struct {
	LocationId    string `uri:"locationId" binding:"required""`
	EventName     string `uri:"eventName" binding:"required" validate:"min=1,max=30"`
	Description   string `uri:"description" binding:"required" validate:"min=1,max=200"`
	ScheduleBegin string `uri:"scheduleBegin" binding:"required"`
	ScheduleEnd   string `uri:"scheduleEnd" binding:"required"`
}

type EventResponse struct {
	ID            string                `json:"id"`
	EventName     string                `json:"eventName""`
	Description   string                `json:"description"`
	ScheduleBegin string                `json:"scheduleBegin"`
	ScheduleEnd   string                `json:"scheduleEnd"`
	Location      LocationResponse      `json:"location"`
	EventTickets  []EventTicketResponse `json:"eventTickets"`
}

type EventTicketRequest struct {
	EventId  string  `uri:"eventId" binding:"required""`
	Type     string  `uri:"type" binding:"required"`
	Quantity uint    `uri:"quantity" binding:"required"`
	Price    float64 `uri:"price" binding:"required"`
}

type EventTicketResponse struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
}
