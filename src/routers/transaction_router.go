package routers

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/interfaces"
	"github.com/arganjava/online-loket/src/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
	"net/http"
)

type TransactionRouter struct {
	eventRepository interfaces.TransactionRepository
}

func NewTransactionRouter(eventRepository interfaces.TransactionRepository) *TransactionRouter {
	return &TransactionRouter{
		eventRepository: eventRepository,
	}
}

func (l TransactionRouter) CreateTransaction(c *gin.Context) {
	var request dto.TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := validator.Validate(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	_, err := l.eventRepository.CreateTransaction(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Transaction created successfully!"})

}

func (l TransactionRouter) GetTransactionInfo(c *gin.Context) {
	id := c.Param("id")
	data, err := l.eventRepository.FindTransactionById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Transaction Get successfully!", "data": buildTransactionResponse(data)})
}

func buildTransactionResponse(data *models.Transaction) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:              data.ID,
		CustomerName:    data.CustomerName,
		CustomerPhone:   data.CustomerPhone,
		CustomerEmail:   data.CustomerEmail,
		TransactionTime: data.TransactionTime,
		OrderQuantity:   data.OrderQuantity,
		EventTicket:     buildEventTicketResponse(data.EventTicket),
	}
}
