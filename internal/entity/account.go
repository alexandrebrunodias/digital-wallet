package entity

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Account struct {
	ID        uuid.UUID
	Customer  *Customer
	Balance   decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newAccount(customer *Customer) (*Account, error) {
	if customer == nil {
		return nil, errors.New("'customer' should not be null")
	}
	now := time.Now()
	return &Account{
		ID:        uuid.New(),
		Customer:  customer,
		Balance:   decimal.Zero,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (a *Account) Credit(amount decimal.Decimal) error {
	if amount.IsNegative() || amount.IsZero() {
		return errors.New("credit a negative or zero 'amount' is not allowed")
	}
	a.Balance = a.Balance.Add(amount)
	return nil
}

func (a *Account) Debit(amount decimal.Decimal) error {
	if amount.IsNegative() {
		return errors.New("debit a negative or zero 'amount' is not allowed")
	}

	if a.Balance.LessThan(amount) {
		return errors.New(
			fmt.Sprintf("insufficient funds | balance: %s - debit amount: %s",
				a.Balance.String(), amount.String()),
		)
	}

	a.Balance = a.Balance.Sub(amount)
	return nil
}
