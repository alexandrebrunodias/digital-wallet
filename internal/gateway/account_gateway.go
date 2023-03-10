package gateway

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type AccountGateway interface {
	Save(account *entity.Account) error
	GetByID(ID uuid.UUID) (*entity.Account, error)
}
