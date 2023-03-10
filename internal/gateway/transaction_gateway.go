package gateway

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type TransactionGateway interface {
	Save(transaction *entity.Transaction) error
	GetByID(ID uuid.UUID) (*entity.Transaction, error)
}
