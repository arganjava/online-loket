package service_repository

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/models"
)

type EventRepository interface {
	CreateEvent(request dto.EventRequest) (int64, error)
	FindEventById(id string) (*models.Event, error)
}
