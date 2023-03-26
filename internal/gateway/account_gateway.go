package gateway

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountGateway interface {
	Create(account *entity.Account) error
	GetByID(ID uuid.UUID) (*entity.Account, error)
	UpdateBalance(ID uuid.UUID, amount decimal.Decimal) error
}
