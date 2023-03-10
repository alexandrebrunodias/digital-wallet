package create_customer

import (
	"errors"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type CustomerGatewayMock struct {
	mock.Mock
}

func (m *CustomerGatewayMock) Save(customer *entity.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *CustomerGatewayMock) GetById(uuid uuid.UUID) (*entity.Customer, error) {
	args := m.Called(uuid)
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func TestCreateCustomerUseCase_CreateSuccessfully(t *testing.T) {
	gatewayMock := &CustomerGatewayMock{}
	gatewayMock.On("Save", mock.Anything).Return(nil)

	expectedName := "alex"
	expectedEmail := "alexandrebrunodias@gmail.com"

	command := CreateCustomerCommand{
		Name:  expectedName,
		Email: expectedEmail,
	}

	useCase := NewCreateCustomerUseCase(gatewayMock)
	output, err := useCase.Execute(command)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, expectedName, output.Name)
	assert.Equal(t, expectedEmail, output.Email)
	assert.NotNil(t, output.CreatedAt)
	assert.NotNil(t, output.UpdatedAt)

	gatewayMock.AssertExpectations(t)
	gatewayMock.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateCustomerUseCase_FailDueToGatewayError(t *testing.T) {
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
