package create_customer

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/alexandrebrunodias/wallet-core/internal/gateway"
	"github.com/google/uuid"
	"time"
)

type CreateCustomerCommand struct {
	Name  string
	Email string
}

type CreateCustomerOutput struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCustomerUseCase struct {
	CustomerGateway gateway.CustomerGateway
}

func NewCreateCustomerUseCase(customerGateway gateway.CustomerGateway) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		CustomerGateway: customerGateway,
	}
}

func (uc *CreateCustomerUseCase) Execute(command CreateCustomerCommand) (*CreateCustomerOutput, error) {
	customer, err := entity.NewCustomer(command.Name, command.Email)
	if err != nil {
		return nil, err
	}

	err = uc.CustomerGateway.Save(customer)

	if err != nil {
		return nil, err
	}

	return &CreateCustomerOutput{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}, nil
}
