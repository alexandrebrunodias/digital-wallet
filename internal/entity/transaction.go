package entity

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Status string

type Transaction struct {
	ID          uuid.UUID
	FromAccount *Account
	ToAccount   *Account
	Status      Status
	Amount      decimal.Decimal
	CreatedAt   time.Time
}

func NewTransaction(fromAccount *Account, toAccount *Account, amount decimal.Decimal) (*Transaction, error) {
	transaction := &Transaction{
		ID:          uuid.New(),
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      amount,
		CreatedAt:   time.Now().UTC(),
	}
	err := transaction.Validate()
	if err != nil {
		return nil, err
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) Commit() error {
	if err := t.FromAccount.Debit(t.Amount); err != nil {
		return err
	}
	if err := t.ToAccount.Credit(t.Amount); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) Validate() error {
	if t.FromAccount == nil || t.ToAccount == nil {
		return errors.New("neither 'FromAccount' nor 'ToAccount' can be nil")
	}
	if t.Amount.IsNegative() || t.Amount.IsZero() {
		return errors.New("'amount' must be a non zero positive number")
	}
	return nil
}
