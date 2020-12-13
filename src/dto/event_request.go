package dto

type EventRequest struct {
	LocationId    string `uri:"locationId" binding:"required""`
	EventName     string `uri:"eventName" binding:"required" validate:"min=1,max=30"`
	Description   string `uri:"description" binding:"required" validate:"min=1,max=200"`
	ScheduleBegin string `uri:"scheduleBegin" binding:"required"`
	ScheduleEnd   string `uri:"scheduleEnd" binding:"required"`
}
