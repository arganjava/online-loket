package service_repository

import (
	"github.com/arganjava/online-loket/src/dto"
	"github.com/arganjava/online-loket/src/models"
)

type TransactionRepository interface {
	CreateTransaction(request dto.TransactionRequest) (int64, error)
	FindTransactionById(id string) (*models.Transaction, error)
}
