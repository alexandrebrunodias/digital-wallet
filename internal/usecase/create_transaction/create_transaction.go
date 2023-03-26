package create_transaction

import (
	"context"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/alexandrebrunodias/wallet-core/internal/gateway"
	"github.com/alexandrebrunodias/wallet-core/pkg/uow"
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
	UnitOfWork uow.UnitOfWorkInterface
}

func NewCreateTransactionUseCase(UnitOfWork uow.UnitOfWorkInterface) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{UnitOfWork}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, command CreateTransactionCommand) (*CreateTransactionOutput, error) {
	output := &CreateTransactionOutput{}
	err := uc.UnitOfWork.Do(ctx, func(unitOfWork *uow.UnitOfWork) error {
		accountGateway := uc.getAccountGateway(ctx)
		transactionGateway := uc.getTransactionGateway(ctx)

		fromAccount, err := accountGateway.GetByID(command.FromAccountID)
		if err != nil {
			return err
		}

		toAccount, err := accountGateway.GetByID(command.ToAccountID)
		if err != nil {
			return err
		}

		transaction, err := entity.NewTransaction(fromAccount, toAccount, command.Amount)
		if err != nil {
			return err
		}

		err = transactionGateway.Save(transaction)
		if err != nil {
			return err
		}

		err = accountGateway.UpdateBalance(fromAccount.ID, transaction.FromAccount.Balance)
		if err != nil {
			return err
		}

		err = accountGateway.UpdateBalance(toAccount.ID, transaction.ToAccount.Balance)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountGateway(ctx context.Context) gateway.AccountGateway {
	repository, err := uc.UnitOfWork.GetRepository(ctx, "AccountGateway")
	if err != nil {
		panic(err)
	}
	return repository.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionGateway(ctx context.Context) gateway.TransactionGateway {
	repository, err := uc.UnitOfWork.GetRepository(ctx, "TransactionGateway")
	if err != nil {
		panic(err)
	}
	return repository.(gateway.TransactionGateway)
}
