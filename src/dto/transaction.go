package dto

import "time"

type TransactionRequest struct {
	EventTicketId string `uri:"eventTicketId" binding:"required"`
	CustomerName  string `uri:"customerName" binding:"required" validate:"min=1,max=100"`
	CustomerPhone string `uri:"customerPhone" binding:"required" validate:"min=1,max=20"`
	CustomerEmail string `uri:"customerEmail" binding:"required" validate:"min=1,max=100"`
	OrderQuantity int64  `uri:"orderQuantity" binding:"required" validate:"min=1"`
}

type TransactionResponse struct {
	ID              string              `json:"id"`
	EventTicket     EventTicketResponse `json:"eventTicket"`
	CustomerName    string              `json:"customerName"`
	CustomerPhone   string              `json:"customerPhone"`
	CustomerEmail   string              `json:"customerEmail"`
	OrderQuantity   int64               `json:"orderQuantity"`
	TransactionTime time.Time           `json:"transactionTime"`
}
