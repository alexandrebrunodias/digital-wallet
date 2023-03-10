package entity

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount_CreateSuccessfully(t *testing.T) {
	expectedCustomer, _ := newCustomer("alex", "alexandrebrunodias@gmail.com")
	expectedBalance := decimal.Zero
	account, err := newAccount(expectedCustomer)

	assert.Nil(t, err)
	assert.Equal(t, expectedCustomer.ID, account.Customer.ID)
	assert.Equal(t, expectedBalance, account.Balance)
}

func TestNewAccount_FailDueToNullCustomer(t *testing.T) {
	expectedErrorMessage := "'customer' should not be null"
	account, err := newAccount(nil)

	assert.NotNil(t, err)
	assert.Nil(t, account)
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestCreditAccount_Successfully(t *testing.T) {
	expectedCustomer, _ := newCustomer("alex", "alexandrebrunodias@gmail.com")
	expectedBalance, _ := decimal.NewFromString("1000.32")
	account, _ := newAccount(expectedCustomer)

	err := account.Credit(expectedBalance)

	assert.Nil(t, err)
	assert.Equal(t, expectedCustomer.ID, account.Customer.ID)
	assert.Equal(t, expectedBalance.String(), account.Balance.String())
}

func TestCreditAccount_FailDueZeroAmount(t *testing.T) {
	expectedCustomer, _ := newCustomer("alex", "alexandrebrunodias@gmail.com")
	expectedErrorMessage := "credit a negative or zero 'amount' is not allowed"
	expectedBalance := decimal.Zero
	account, _ := newAccount(expectedCustomer)

	err := account.Credit(expectedBalance)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestCreditAccount_FailDueNegativeAmount(t *testing.T) {
	expectedCustomer, _ := newCustomer("alex", "alexandrebrunodias@gmail.com")
	expectedErrorMessage := "credit a negative or zero 'amount' is not allowed"
	expectedBalance, _ := decimal.NewFromString("-1000.32")
	account, _ := newAccount(expectedCustomer)

	err := account.Credit(expectedBalance)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestDebitAccount_Successfully(t *testing.T) {
	expectedCustomer, _ := newCustomer("alex", "alexandrebrunodias@gmail.com")
	expectedBalance, _ := decimal.NewFromString("1000.32")
	account, _ := newAccount(expectedCustomer)

	err := account.Credit(expectedBalance)

	assert.Nil(t, err)
	assert.Equal(t, expectedCustomer.ID, account.Customer.ID)
	assert.Equal(t, expectedBalance.String(), account.Balance.String())
}

func TestDebitAccount_FailDueNegativeAmount(t *testing.T) {
	expectedCustomer, _ := newCustomer("alex", "alexandrebrunodias@gmail.com")
	account, _ := newAccount(expectedCustomer)
	expectedErrorMessage := "debit a negative or zero 'amount' is not allowed"

	negativeAmount, _ := decimal.NewFromString("-1000.32")
	err := account.Debit(negativeAmount)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestDebitAccount_FailDueToInsufficientFunds(t *testing.T) {
	expectedCustomer, _ := newCustomer("alex", "alexandrebrunodias@gmail.com")
	account, _ := newAccount(expectedCustomer)

	amount, _ := decimal.NewFromString("1000.32")
	expectedErrorMessage := "insufficient funds | balance: 0 - debit amount: " + amount.String()

	err := account.Debit(amount)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
}
