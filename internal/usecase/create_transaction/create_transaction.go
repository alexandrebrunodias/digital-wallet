package create_transaction

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/alexandrebrunodias/wallet-core/internal/gateway"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateTransactionCommand struct {
	FromAccountID uuid.UUID
	ToAccountID   uuid.UUID
	Amount        decimal.Decimal
}

type CreateTransactionOutput struct {
	ID uuid.UUID
}

type CreateTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway     gateway.AccountGateway
}

func NewCreateTransactionUseCase(
	transactionGateway gateway.TransactionGateway,
	accountGateway gateway.AccountGateway,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: transactionGateway,
		AccountGateway:     accountGateway,
	}
}

func (uc *CreateTransactionUseCase) Execute(command CreateTransactionCommand) (*CreateTransactionOutput, error) {
	fromAccount, err := uc.AccountGateway.GetByID(command.FromAccountID)
	if err != nil {
		return nil, err
	}
	toAccount, err := uc.AccountGateway.GetByID(command.ToAccountID)
	if err != nil {
		return nil, err
	}

	transaction, err := entity.NewTransaction(fromAccount, toAccount, command.Amount)
	if err != nil {
		return nil, err
	}
	err = uc.TransactionGateway.Save(transaction)
	if err != nil {
		return nil, err
	}

	return &CreateTransactionOutput{ID: transaction.ID}, nil
}
