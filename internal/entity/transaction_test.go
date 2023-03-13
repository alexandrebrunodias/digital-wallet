package entity

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestNewTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}

func (s *TransactionTestSuite) TestNewTransaction_CreateSuccessfully() {
	expectedStatus := COMPLETED
	expectedAmount := decimal.NewFromInt(100)

	expectedFromAccountInitialBalance := decimal.NewFromInt(200)
	expectedFromAccountFinalBalance := expectedFromAccountInitialBalance.Sub(expectedAmount)
	_ = s.AccountFrom.Credit(expectedFromAccountInitialBalance)

	expectedToAccountInitialBalance := decimal.NewFromInt(200)
	expectedToAccountFinalBalance := expectedToAccountInitialBalance.Add(expectedAmount)
	_ = s.AccountTo.Credit(expectedToAccountInitialBalance)

	transaction, err := NewTransaction(s.AccountFrom, s.AccountTo, expectedAmount)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), transaction)
	assert.Equal(s.T(), expectedStatus, transaction.Status)
	assert.Equal(s.T(), expectedAmount, transaction.Amount)
	assert.Equal(s.T(), expectedFromAccountFinalBalance, transaction.fromAccount.Balance)
	assert.Equal(s.T(), expectedToAccountFinalBalance, transaction.toAccount.Balance)
	assert.NotNil(s.T(), transaction.CreatedAt)
}

func (s *TransactionTestSuite) TestNewTransaction_FailDueToNilFromAccount() {
	expectedErrorMessage := "neither 'fromAccount' nor 'toAccount' can be nil"

	transaction, err := NewTransaction(nil, s.AccountTo, decimal.NewFromInt(100))

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErrorMessage, err.Error())
	assert.Nil(s.T(), transaction)
}

func (s *TransactionTestSuite) TestNewTransaction_FailDueToNilToAccount() {
	expectedErrorMessage := "neither 'fromAccount' nor 'toAccount' can be nil"

	transaction, err := NewTransaction(nil, s.AccountFrom, decimal.NewFromInt(100))

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErrorMessage, err.Error())
	assert.Nil(s.T(), transaction)
}

func (s *TransactionTestSuite) TestNewTransaction_FailDueToNegativeAmount() {
	expectedErrorMessage := "'amount' must be a non zero positive number"

	transaction, err := NewTransaction(s.AccountFrom, s.AccountTo, decimal.NewFromInt(-100))

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErrorMessage, err.Error())
	assert.Nil(s.T(), transaction)
}

func (s *TransactionTestSuite) TestNewTransaction_FailDueToNegativeZero() {
	expectedErrorMessage := "'amount' must be a non zero positive number"

	transaction, err := NewTransaction(s.AccountFrom, s.AccountTo, decimal.Zero)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErrorMessage, err.Error())
	assert.Nil(s.T(), transaction)
}

func (s *TransactionTestSuite) TestNewTransaction_FailDueToAlreadyCompleted() {
	_ = s.AccountFrom.Credit(decimal.NewFromInt(200))

	expectedStatus := COMPLETED
	expectedAmount := decimal.NewFromInt(100)
	expectedErrorMessage := "transaction is already COMPLETED"

	transaction, _ := NewTransaction(s.AccountFrom, s.AccountTo, expectedAmount)

	err := transaction.Commit()

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErrorMessage, err.Error())
	assert.NotNil(s.T(), transaction)
	assert.Equal(s.T(), expectedStatus, transaction.Status)
	assert.Equal(s.T(), expectedAmount, transaction.Amount)
	assert.Equal(s.T(), s.AccountFrom, transaction.fromAccount)
	assert.Equal(s.T(), s.AccountTo.Balance.String(), expectedAmount.String())
	assert.NotNil(s.T(), transaction.CreatedAt)
}

func (s *TransactionTestSuite) TestNewTransaction_FailDueToErrorDebitingAccount() {
	expectedAmount := decimal.NewFromInt(100)
	expectedErrorMessage := fmt.Sprintf(
		"customer %s has insufficient funds | balance: 0 - debit amount: %s",
		s.CustomerFrom.ID, expectedAmount.String(),
	)

	transaction, err := NewTransaction(s.AccountFrom, s.AccountTo, expectedAmount)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedErrorMessage, err.Error())
	assert.Nil(s.T(), transaction)
}

type TransactionTestSuite struct {
	suite.Suite
	CustomerFrom *Customer
	CustomerTo   *Customer
	AccountFrom  *Account
	AccountTo    *Account
}

func (s *TransactionTestSuite) SetupTest() {
	var err error
	s.CustomerFrom, err = NewCustomer("alex", "alexandrebrunodias@gmail.com")
	s.Require().Nil(err)
	s.AccountFrom, err = NewAccount(s.CustomerFrom)
	s.Require().Nil(err)

	s.CustomerTo, err = NewCustomer("alex", "alexandrebrunodias@gmail.com")
	s.Require().Nil(err)
	s.AccountTo, err = NewAccount(s.CustomerTo)
	s.Require().Nil(err)
}
