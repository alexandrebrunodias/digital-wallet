package create_customer

import (
	"errors"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateCustomerUseCase_Execute_CreateSuccessfully(t *testing.T) {
	customerGatewayMock := &CustomerGatewayMock{}
	customerGatewayMock.On("Save", mock.Anything).Return(nil)

	expectedName := "alex"
	expectedEmail := "alexandrebrunodias@gmail.com"

	command := CreateCustomerCommand{
		Name:  expectedName,
		Email: expectedEmail,
	}

	useCase := NewCreateCustomerUseCase(customerGatewayMock)
	output, err := useCase.Execute(command)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, expectedName, output.Name)
	assert.Equal(t, expectedEmail, output.Email)
	assert.NotNil(t, output.CreatedAt)
	assert.NotNil(t, output.UpdatedAt)

	customerGatewayMock.AssertExpectations(t)
	customerGatewayMock.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateCustomerUseCase_Execute_FailDueToGatewayError(t *testing.T) {
	expectedErrorMessage := "gateway error"

	gatewayMock := &CustomerGatewayMock{}
	gatewayMock.On("Save", mock.Anything).Return(errors.New(expectedErrorMessage))

	expectedName := "alex"
	expectedEmail := "alexandrebrunodias@gmail.com"

	command := CreateCustomerCommand{
		Name:  expectedName,
		Email: expectedEmail,
	}

	useCase := NewCreateCustomerUseCase(gatewayMock)
	output, err := useCase.Execute(command)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErrorMessage, err.Error())
	assert.Nil(t, output)

	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Save", 1)
}

type CustomerGatewayMock struct {
	mock.Mock
}

func (m *CustomerGatewayMock) Save(customer *entity.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *CustomerGatewayMock) GetByID(ID uuid.UUID) (*entity.Customer, error) {
	args := m.Called(ID)
	return args.Get(0).(*entity.Customer), args.Error(1)
}
