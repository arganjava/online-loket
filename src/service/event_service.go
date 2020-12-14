package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/interfaces"
	"github.com/arganjava/online-loket/src/models"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type EventService struct {
	db           *sql.DB
	locationRepo interfaces.LocationRepository
}

func NewEventService(db *sql.DB, locationRepository interfaces.LocationRepository) *EventService {
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
		return 0, fmt.Errorf("Event already exist for %v %v from %v to %v ", location.CityName, request.EventName, request.ScheduleBegin, request.ScheduleEnd)
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

func (l EventService) CreateEventTicket(request dto.EventTicketRequest) (int64, error) {
	event, err := l.FindEventById(request.EventId)
	if err != nil {
		return 0, err
	}

	if event == nil {
		return 0, fmt.Errorf("Event not found for id %v", request.EventId)
	}

	isExist, err := l.isEventTicketExist(request)
	if err != nil {
		return 0, err
	}

	if isExist {
		return 0, fmt.Errorf("Event Ticket already exist for %v %v ", event.Name, request.Type)
	}

	ctx := context.Background()
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	id := uuid.NewV4()
	sid := id.String()
	sql := fmt.Sprintf("INSERT INTO event_ticket (id, event_id, ticket_type, quantity, price) VALUES ('%v',  '%v',  '%v',  %v,  %v)",
		sid, request.EventId, request.Type, request.Quantity, request.Price)
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

func (l EventService) FindEventTicketId(id string) (*models.EventTicket, error) {
	rows, err := l.db.Query("SELECT id, ticket_type, quantity, price, event_id FROM event_ticket "+
		"WHERE id = $1",
		id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		ticket := &models.EventTicket{}
		var eventId string
		err = rows.Scan(&ticket.ID, &ticket.Type, &ticket.Quantity, &ticket.Price, &eventId)
		if err != nil {
			return nil, err
		} else {
			return ticket, nil
		}
	}

	return nil, nil
}

func (l EventService) FindEventById(id string) (*models.Event, error) {
	l.db.Begin()
	rows, err := l.db.Query("SELECT id, event_name, description, schedule_begin, schedule_end, location_id  FROM event "+
		"WHERE id = $1",
		id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		event := &models.Event{}
		var locationId string
		err = rows.Scan(&event.ID, &event.Name, &event.Description, &event.ScheduleBegin, &event.ScheduleEnd, &locationId)
		if err != nil {
			return nil, err
		} else {
			location, err := l.locationRepo.FindLocationById(locationId)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			eventTickets, err := l.findTicketsByEventId(event.ID)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			event.Location = location
			event.EventTickets = eventTickets
			return event, nil
		}
	}
	return nil, nil
}

func (l EventService) findTicketsByEventId(id string) ([]*models.EventTicket, error) {
	rows, err := l.db.Query("SELECT id, ticket_type, price, quantity  FROM event_ticket "+
		"WHERE event_id = $1",
		id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	evetTickets := make([]*models.EventTicket, 0)
	for rows.Next() {
		if rows.Next() {
			eventTicket := &models.EventTicket{}
			err = rows.Scan(&eventTicket.ID, &eventTicket.Type, &eventTicket.Price, &eventTicket.Quantity)
			if err != nil {
				return nil, err
			} else {
				evetTickets = append(evetTickets, eventTicket)
			}
		}
	}

	return evetTickets, nil
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
	defer rows.Close()
	if err != nil {
		log.Print(err)
		return false, err
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (l EventService) isEventTicketExist(request dto.EventTicketRequest) (bool, error) {
	rows, err := l.db.Query("SELECT event_id, ticket_type  FROM event_ticket "+
		"WHERE  event_id= $1 and ticket_type = $2",
		request.EventId, request.Type)
	if err != nil {
		log.Print(err)
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}
