package create_account

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/alexandrebrunodias/wallet-core/internal/gateway"
	"github.com/google/uuid"
)

type CreateAccountCommand struct {
	CustomerID uuid.UUID
}

type CreateAccountOutput struct {
	ID uuid.UUID
}

type CreateAccountUseCase struct {
	AccountGateway  gateway.AccountGateway
	CustomerGateway gateway.CustomerGateway
}

func NewAccountUseCase(
	accountGateway gateway.AccountGateway,
	customerGateway gateway.CustomerGateway,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway:  accountGateway,
		CustomerGateway: customerGateway,
	}
}

func (uc *CreateAccountUseCase) Execute(command CreateAccountCommand) (*CreateAccountOutput, error) {
	customer, err := uc.CustomerGateway.GetByID(command.CustomerID)
	if err != nil {
		return nil, err
	}

	account, err := entity.NewAccount(customer)
	if err != nil {
		return nil, err
	}

	err = uc.AccountGateway.Create(account)
	if err != nil {
		return nil, err
	}

	return &CreateAccountOutput{ID: account.ID}, nil
}
