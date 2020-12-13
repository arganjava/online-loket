package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/models"
	"github.com/arganjava/online-loket/src/service_repository"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type EventService struct {
	db           *sql.DB
	locationRepo service_repository.LocationRepository
}

func NewEventService(db *sql.DB, locationRepository service_repository.LocationRepository) *EventService {
	return &EventService{
		db:           db,
		locationRepo: locationRepository,
	}
}

func (l EventService) CreateEvent(request dto.EventRequest) (int64, error) {
	location, err := l.locationRepo.FindLocationById(request.LocationId)
	if err != nil {
		return 0, err
	}

	if location == nil {
		return 0, fmt.Errorf("Location not found for id %v", request.LocationId)
	}

	isExist, err := l.isEventExist(request)
	if err != nil {
		return 0, err
	}

	if isExist {
		return 0, fmt.Errorf("Event already exist for %v %v from %v to %v ", location.CityName.String, request.EventName, request.ScheduleBegin, request.ScheduleEnd)
	}

	ctx := context.Background()
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	id := uuid.NewV4()
	sid := id.String()

	sql := fmt.Sprintf("INSERT INTO event (id, event_name, description, schedule_begin, schedule_end, location_id) VALUES ('%v',  '%v',  '%v',  '%v',  '%v',  '%v')",
		sid, request.EventName, request.Description, request.ScheduleBegin, request.ScheduleEnd, request.LocationId)
	result, err := tx.ExecContext(ctx, sql)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return result.RowsAffected()
}

func (l EventService) FindEventById(id string) (*models.Event, error) {
	rows, err := l.db.Query("SELECT id, event_name, description, schedule_begin, schedule_end, location_id  FROM event "+
		"WHERE id = $1",
		id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if rows.Next() {
		var id string
		var eventName string
		var description string
		var scheduleBegin time.Time
		var scheduleEnd time.Time
		var locationId string
		err = rows.Scan(&id, &eventName, &description, &scheduleBegin, &scheduleEnd, &locationId)
		if err != nil {
			return nil, err
		} else {
			location := models.Location{ID: locationId}
			return &models.Event{id, location, eventName, description, scheduleBegin, scheduleEnd}, nil
		}
	}
	return nil, nil
}

func (l EventService) isEventExist(request dto.EventRequest) (bool, error) {
	scheduleBegin, err := time.Parse("2006-01-02", request.ScheduleBegin)
	if err != nil {
		return false, err
	}
	scheduleEnd, err := time.Parse("2006-01-02", request.ScheduleEnd)
	if err != nil {
		return false, err
	}

	rows, err := l.db.Query("SELECT event_name, schedule_begin, schedule_end, location_id  FROM event "+
		"WHERE  location_id= $1 and event_name = $2 and schedule_begin = $3 and schedule_end = $4",
		request.LocationId, request.EventName,
		scheduleBegin, scheduleEnd)
	if err != nil {
		log.Print(err)
		return false, err
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
}
