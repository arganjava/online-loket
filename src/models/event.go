package models

import "time"

type Event struct {
	ID            string
	Location      Location
	Name          string
	Description   string
	ScheduleBegin time.Time
	ScheduleEnd   time.Time
}
