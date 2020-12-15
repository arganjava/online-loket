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
)

type TransactionService struct {
	db        *sql.DB
	eventRepo interfaces.EventRepository
}

func NewTransactionService(db *sql.DB, eventRepo interfaces.EventRepository) *TransactionService {
	return &TransactionService{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (l TransactionService) CreateTransaction(request dto.TransactionRequest) (int64, error) {

	ctx := context.Background()
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	ticket, err := l.eventRepo.FindEventTicketId(request.EventTicketId)
	if err != nil {
		return 0, err
	}

	if ticket == nil {
		return 0, fmt.Errorf("Event Ticket not found for id %v ", request.EventTicketId)

	}

	isAvailable, remainQuantity := l.isAvailableQuota(ticket, request)

	if !isAvailable {
		return 0, fmt.Errorf("Event Ticket not enough quota for %v ", ticket.ID)
	}

	id := uuid.NewV4()
	sid := id.String()

	sqlInsertTrx := fmt.Sprintf("INSERT INTO transaction (id, customer_name, customer_phone, customer_email, order_quantity, transaction_time, event_ticket_id) VALUES ('%v',  '%v',  '%v',  '%v',  '%v',  NOW(), '%v')",
		sid, request.CustomerName, request.CustomerPhone, request.CustomerEmail, request.OrderQuantity, request.EventTicketId)
	result, err := tx.ExecContext(ctx, sqlInsertTrx)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return 0, err
	}

	sqlUpdateTicket := fmt.Sprintf("UPDATE event_ticket SET quantity= %v WHERE id = '%v'", remainQuantity, request.EventTicketId)
	result, err = tx.ExecContext(ctx, sqlUpdateTicket)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return 0, err
	}

	return result.RowsAffected()
}

func (l TransactionService) isAvailableQuota(ticket *models.EventTicket,
	request dto.TransactionRequest) (bool, int64) {

	reduceResult := ticket.Quantity - request.OrderQuantity
	if reduceResult < 1 {
		return false, reduceResult
	}

	return true, reduceResult

}

func (l TransactionService) FindTransactionById(id string) (*models.Transaction, error) {
	rows, err := l.db.Query("SELECT id, customer_name, customer_phone, customer_email, order_quantity, transaction_time, event_ticket_id  FROM transaction "+
		"WHERE id = $1",
		id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		transaction := &models.Transaction{}
		var eventTicketId string
		err = rows.Scan(&transaction.ID, &transaction.CustomerName, &transaction.CustomerPhone, &transaction.CustomerEmail, &transaction.OrderQuantity, &transaction.TransactionTime, &eventTicketId)
		if err != nil {
			return nil, err
		} else {
			ticket, err := l.eventRepo.FindEventTicketId(eventTicketId)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			transaction.EventTicket = ticket
			return transaction, nil
		}
	}

	return nil, nil
}
