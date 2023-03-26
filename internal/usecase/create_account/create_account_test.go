package create_account

import (
	"errors"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateAccountUseCase_Execute_CreateSuccessfully(t *testing.T) {
	customer, _ := entity.NewCustomer("alex", "alexandrebrunodias@gmail.com")
	customerGatewayMock := &CustomerGatewayMock{}
	accountGatewayMock := &AccountGatewayMock{}

	customerGatewayMock.On("GetByID", customer.ID).
		Return(customer, nil)

	accountGatewayMock.On("Create", m.AnythingOfType("*entity.Account")).
		Return(nil)

	command := CreateAccountCommand{customer.ID}

	useCase := NewAccountUseCase(accountGatewayMock, customerGatewayMock)
	output, err := useCase.Execute(command)

	assert.Nil(t, err)
	assert.NotNil(t, output.ID)

	customerGatewayMock.AssertExpectations(t)
	customerGatewayMock.AssertNumberOfCalls(t, "GetByID", 1)

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "Create", 1)
}

func TestCreateAccountUseCase_Execute_FailDueToCustomerNotFound(t *testing.T) {
	customerID := uuid.New()
	expectedErrorMessage := "customer not found"

	customerGatewayMock := &CustomerGatewayMock{}
	accountGatewayMock := &AccountGatewayMock{}

	customerGatewayMock.On("GetByID", customerID).
		Return(&entity.Customer{}, errors.New(expectedErrorMessage))

	command := CreateAccountCommand{customerID}

	useCase := NewAccountUseCase(accountGatewayMock, customerGatewayMock)
	output, err := useCase.Execute(command)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, output)

	customerGatewayMock.AssertExpectations(t)
	customerGatewayMock.AssertNumberOfCalls(t, "GetByID", 1)

	accountGatewayMock.AssertNotCalled(t, "Create")
}

func TestCreateAccountUseCase_Execute_FailDueToGatewayError(t *testing.T) {
	customer, _ := entity.NewCustomer("alex", "alexandrebrunodias@gmail.com")
	customerGatewayMock := &CustomerGatewayMock{}
	accountGatewayMock := &AccountGatewayMock{}
	expectedErrorMessage := "gateway error"

	customerGatewayMock.On("GetByID", customer.ID).
		Return(customer, nil)

	accountGatewayMock.On("Create", m.AnythingOfType("*entity.Account")).
		Return(errors.New(expectedErrorMessage))

	command := CreateAccountCommand{customer.ID}

	useCase := NewAccountUseCase(accountGatewayMock, customerGatewayMock)
	output, err := useCase.Execute(command)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, output)

	customerGatewayMock.AssertExpectations(t)
	customerGatewayMock.AssertNumberOfCalls(t, "GetByID", 1)

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "Create", 1)
}

type AccountGatewayMock struct {
	m.Mock
}

func (m *AccountGatewayMock) UpdateBalance(ID uuid.UUID, amount decimal.Decimal) error {
	//TODO implement me
	panic("implement me")
}

func (m *AccountGatewayMock) Create(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) GetByID(ID uuid.UUID) (*entity.Account, error) {
	args := m.Called(ID)
	return args.Get(0).(*entity.Account), args.Error(1)
}

type CustomerGatewayMock struct {
	m.Mock
}

func (m *CustomerGatewayMock) Create(customer *entity.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *CustomerGatewayMock) GetByID(ID uuid.UUID) (*entity.Customer, error) {
	args := m.Called(ID)
	return args.Get(0).(*entity.Customer), args.Error(1)
}
