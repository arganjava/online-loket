package models

import (
	"time"
)

type Event struct {
	ID            string
	Location      *Location
	Name          string
	Description   string
	ScheduleBegin time.Time
	ScheduleEnd   time.Time
	EventTickets  []*EventTicket
}

type EventTicket struct {
	ID       string
	Type     string
	Quantity int64
	Price    float64
}
