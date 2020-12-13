package models

import (
	"time"
)

type Transaction struct {
	ID              string
	EventTicket     *EventTicket
	CustomerName    string
	CustomerPhone   string
	CustomerEmail   string
	OrderQuantity   int64
	TransactionTime time.Time
}
