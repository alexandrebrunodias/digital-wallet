package entity

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTransaction_CreateSuccessfully(t *testing.T) {
	customer1, _ := NewCustomer("alex", "alexandrebrunodias@gmail.com")
	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")

	expectedFromAccount, _ := newAccount(customer1)
	expectedFromAccount.Credit(decimal.NewFromInt(200))

	expectedToAccount, _ := newAccount(customer2)

	expectedStatus := COMPLETED
	expectedAmount := decimal.NewFromInt(100)

	transaction, err := newTransaction(expectedFromAccount, expectedToAccount, expectedAmount)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, expectedStatus, transaction.Status)
	assert.Equal(t, expectedAmount, transaction.Amount)
	assert.Equal(t, expectedFromAccount, transaction.fromAccount)
	assert.Equal(t, expectedToAccount.Balance, expectedAmount)
}

func TestNewTransaction_FailDueToNilFromAccount(t *testing.T) {
	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := newAccount(customer2)

	expectedErrorMessage := "neither 'fromAccount' nor 'toAccount' can be nil"

	transaction, err := newTransaction(nil, expectedToAccount, decimal.NewFromInt(100))

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToNilToAccount(t *testing.T) {
	customer, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := newAccount(customer)

	expectedErrorMessage := "neither 'fromAccount' nor 'toAccount' can be nil"

	transaction, err := newTransaction(nil, expectedFromAccount, decimal.NewFromInt(100))

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToNegativeAmount(t *testing.T) {
	customer1, _ := NewCustomer("alex1", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := newAccount(customer1)

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := newAccount(customer2)

	expectedErrorMessage := "'amount' must be a non zero positive number"

	transaction, err := newTransaction(expectedFromAccount, expectedToAccount, decimal.NewFromInt(-100))

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToNegativeZero(t *testing.T) {
	customer1, _ := NewCustomer("alex1", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := newAccount(customer1)

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := newAccount(customer2)

	expectedErrorMessage := "'amount' must be a non zero positive number"

	transaction, err := newTransaction(expectedFromAccount, expectedToAccount, decimal.Zero)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToAlreadyCompleted(t *testing.T) {
	customer1, _ := NewCustomer("alex1", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := newAccount(customer1)
	expectedFromAccount.Credit(decimal.NewFromInt(200))

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := newAccount(customer2)

	expectedStatus := COMPLETED
	expectedAmount := decimal.NewFromInt(100)
	expectedErrorMessage := "transaction is already COMPLETED"

	transaction, _ := newTransaction(expectedFromAccount, expectedToAccount, expectedAmount)

	err := transaction.Commit()

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.NotNil(t, transaction)
	assert.Equal(t, expectedStatus, transaction.Status)
	assert.Equal(t, expectedAmount, transaction.Amount)
	assert.Equal(t, expectedFromAccount, transaction.fromAccount)
	assert.Equal(t, expectedToAccount.Balance, expectedAmount)
}

func TestNewTransaction_FailDueToErrorDebitingAccount(t *testing.T) {
	customer1, _ := NewCustomer("alex1", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := newAccount(customer1)

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := newAccount(customer2)

	expectedAmount := decimal.NewFromInt(100)
	expectedErrorMessage := "insufficient funds | balance: 0 - debit amount: " + expectedAmount.String()

	transaction, err := newTransaction(expectedFromAccount, expectedToAccount, expectedAmount)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}
