package entity

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTransaction_CreateSuccessfully(t *testing.T) {
	customer1, _ := NewCustomer("alex", "alexandrebrunodias@gmail.com")
	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")

	expectedStatus := COMPLETED
	expectedAmount := decimal.NewFromInt(100)

	expectedFromAccount, _ := NewAccount(customer1)
	expectedFromAccountInitialBalance := decimal.NewFromInt(200)
	expectedFromAccountFinalBalance := expectedFromAccountInitialBalance.Sub(expectedAmount)
	expectedFromAccount.Credit(expectedFromAccountInitialBalance)

	expectedToAccount, _ := NewAccount(customer2)
	expectedToAccountInitialBalance := decimal.NewFromInt(200)
	expectedToAccountFinalBalance := expectedToAccountInitialBalance.Add(expectedAmount)
	expectedToAccount.Credit(expectedToAccountInitialBalance)

	transaction, err := NewTransaction(expectedFromAccount, expectedToAccount, expectedAmount)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, expectedStatus, transaction.Status)
	assert.Equal(t, expectedAmount, transaction.Amount)
	assert.Equal(t, expectedFromAccountFinalBalance, transaction.fromAccount.Balance)
	assert.Equal(t, expectedToAccountFinalBalance, transaction.toAccount.Balance)
}

func TestNewTransaction_FailDueToNilFromAccount(t *testing.T) {
	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := NewAccount(customer2)

	expectedErrorMessage := "neither 'fromAccount' nor 'toAccount' can be nil"

	transaction, err := NewTransaction(nil, expectedToAccount, decimal.NewFromInt(100))

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToNilToAccount(t *testing.T) {
	customer, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := NewAccount(customer)

	expectedErrorMessage := "neither 'fromAccount' nor 'toAccount' can be nil"

	transaction, err := NewTransaction(nil, expectedFromAccount, decimal.NewFromInt(100))

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToNegativeAmount(t *testing.T) {
	customer1, _ := NewCustomer("alex1", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := NewAccount(customer1)

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := NewAccount(customer2)

	expectedErrorMessage := "'amount' must be a non zero positive number"

	transaction, err := NewTransaction(expectedFromAccount, expectedToAccount, decimal.NewFromInt(-100))

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToNegativeZero(t *testing.T) {
	customer1, _ := NewCustomer("alex1", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := NewAccount(customer1)

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := NewAccount(customer2)

	expectedErrorMessage := "'amount' must be a non zero positive number"

	transaction, err := NewTransaction(expectedFromAccount, expectedToAccount, decimal.Zero)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}

func TestNewTransaction_FailDueToAlreadyCompleted(t *testing.T) {
	customer1, _ := NewCustomer("alex1", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := NewAccount(customer1)
	expectedFromAccount.Credit(decimal.NewFromInt(200))

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := NewAccount(customer2)

	expectedStatus := COMPLETED
	expectedAmount := decimal.NewFromInt(100)
	expectedErrorMessage := "transaction is already COMPLETED"

	transaction, _ := NewTransaction(expectedFromAccount, expectedToAccount, expectedAmount)

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
	expectedFromAccount, _ := NewAccount(customer1)

	customer2, _ := NewCustomer("alex2", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := NewAccount(customer2)

	expectedAmount := decimal.NewFromInt(100)
	expectedErrorMessage := fmt.Sprintf(
		"customer %s has insufficient funds | balance: 0 - debit amount: %s",
		customer1.ID, expectedAmount.String(),
	)

	transaction, err := NewTransaction(expectedFromAccount, expectedToAccount, expectedAmount)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, transaction)
}
