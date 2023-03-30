package create_transaction

import (
	"context"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/alexandrebrunodias/wallet-core/internal/gateway"
	"github.com/alexandrebrunodias/wallet-core/pkg/events"
	"github.com/alexandrebrunodias/wallet-core/pkg/uow"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const TransactionCreated = "wallet.core.transaction.created"

type CreateTransactionCommand struct {
	FromAccountID uuid.UUID       `json:"from_account_id"`
	ToAccountID   uuid.UUID       `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
}

type CreateTransactionOutput struct {
	ID            uuid.UUID       `json:"id"`
	FromAccountID uuid.UUID       `json:"from_account_id"`
	ToAccountID   uuid.UUID       `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
}

type CreateTransactionUseCase struct {
	UnitOfWork     uow.UnitOfWorkInterface
	EventPublisher events.EventPublisherInterface
}

func NewCreateTransactionUseCase(
	UnitOfWork uow.UnitOfWorkInterface,
	EventPublisher events.EventPublisherInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		UnitOfWork:     UnitOfWork,
		EventPublisher: EventPublisher,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, command CreateTransactionCommand) (*CreateTransactionOutput, error) {
	output := &CreateTransactionOutput{}
	err := uc.UnitOfWork.Do(ctx, func(_ *uow.UnitOfWork) error {
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

		err = accountGateway.UpdateBalance(fromAccount.ID, transaction.FromAccount.Balance)
		if err != nil {
			return err
		}

		err = accountGateway.UpdateBalance(toAccount.ID, transaction.ToAccount.Balance)
		if err != nil {
			return err
		}

		err = transactionGateway.Create(transaction)
		if err != nil {
			return err
		}

		output.ID = transaction.ID
		output.FromAccountID = fromAccount.ID
		output.ToAccountID = toAccount.ID
		output.Amount = transaction.Amount
		return nil
	})
	if err != nil {
		return nil, err
	}

	event := events.NewEvent(TransactionCreated, output)
	uc.EventPublisher.Register(*event).Publish()
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
