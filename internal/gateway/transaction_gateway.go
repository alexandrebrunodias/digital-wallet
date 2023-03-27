package gateway

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
	GetByID(ID uuid.UUID) (*entity.Transaction, error)
}
