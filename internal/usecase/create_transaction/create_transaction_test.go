package create_transaction

import (
	"errors"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateTransactionUseCase_Execute_CreateSuccessfully(t *testing.T) {
	fromCustomer, _ := entity.NewCustomer("fromCustomer", "alexandrebrunodias@gmail.com")
	expectedFromAccount, _ := entity.NewAccount(fromCustomer)
	expectedFromAccount.Credit(decimal.NewFromInt(2000))

	toCustomer, _ := entity.NewCustomer("toCustomer", "alexandrebrunodias@gmail.com")
	expectedToAccount, _ := entity.NewAccount(toCustomer)

	expectedAmount := decimal.NewFromInt(1000)

	command := CreateTransactionCommand{
		FromAccountID: expectedFromAccount.ID,
		ToAccountID:   expectedToAccount.ID,
		Amount:        expectedAmount,
	}

	accountGatewayMock := &AccountGatewayMock{}
	transactionGatewayMock := &TransactionGatewayMock{}

	accountGatewayMock.On("GetByID", expectedFromAccount.ID).
		Return(expectedFromAccount, nil)
	accountGatewayMock.On("GetByID", expectedToAccount.ID).
		Return(expectedToAccount, nil)

	transactionGatewayMock.On("Save", m.AnythingOfType("*entity.Transaction")).
		Return(nil)

	useCase := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	output, err := useCase.Execute(command)

	assert.Nil(t, err)
	assert.NotNil(t, output.ID)

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "GetByID", 2)

	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateTransactionUseCase_Execute_FailDueToFromAccountNotFound(t *testing.T) {
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	expectedErrorMessage := "from account not found"

	command := CreateTransactionCommand{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        decimal.NewFromInt(1000),
	}

	accountGatewayMock := &AccountGatewayMock{}
	transactionGatewayMock := &TransactionGatewayMock{}

	accountGatewayMock.On("GetByID", fromAccountID).
		Return(&entity.Account{}, errors.New(expectedErrorMessage))

	useCase := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	output, err := useCase.Execute(command)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErrorMessage, err.Error())

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "GetByID", 1)

	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNotCalled(t, "Save")
}

func TestCreateTransactionUseCase_Execute_FailDueToInsufficientFunds(t *testing.T) {
	fromAccountID := uuid.New()
	toAccountID := uuid.New()

	expectedErrorMessage := "to account not found"

	expectedAmount := decimal.NewFromInt(1000)

	command := CreateTransactionCommand{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        expectedAmount,
	}

	accountGatewayMock := &AccountGatewayMock{}
	transactionGatewayMock := &TransactionGatewayMock{}

	accountGatewayMock.On("GetByID", fromAccountID).
		Return(&entity.Account{}, nil)
	accountGatewayMock.On("GetByID", toAccountID).
		Return(&entity.Account{}, errors.New(expectedErrorMessage))

	useCase := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	output, err := useCase.Execute(command)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErrorMessage, err.Error())

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "GetByID", 2)

	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNotCalled(t, "Save")
}

func TestCreateTransactionUseCase_Execute_FailDueToAccountToNotFound(t *testing.T) {
	fromAccountID := uuid.New()
	toAccountID := uuid.New()

	expectedErrorMessage := "to account not found"

	expectedAmount := decimal.NewFromInt(1000)

	command := CreateTransactionCommand{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        expectedAmount,
	}

	accountGatewayMock := &AccountGatewayMock{}
	transactionGatewayMock := &TransactionGatewayMock{}

	accountGatewayMock.On("GetByID", fromAccountID).
		Return(&entity.Account{}, nil)
	accountGatewayMock.On("GetByID", toAccountID).
		Return(&entity.Account{}, errors.New(expectedErrorMessage))

	useCase := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)
	output, err := useCase.Execute(command)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, expectedErrorMessage, err.Error())

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "GetByID", 2)

	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNotCalled(t, "Save")
}

type TransactionGatewayMock struct {
	m.Mock
}

func (m *TransactionGatewayMock) Save(customer *entity.Transaction) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *TransactionGatewayMock) GetByID(ID uuid.UUID) (*entity.Transaction, error) {
	args := m.Called(ID)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

type AccountGatewayMock struct {
	m.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) GetByID(ID uuid.UUID) (*entity.Account, error) {
	args := m.Called(ID)
	return args.Get(0).(*entity.Account), args.Error(1)
}
