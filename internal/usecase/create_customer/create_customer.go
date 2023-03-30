package create_customer

import (
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/alexandrebrunodias/wallet-core/internal/gateway"
	"github.com/google/uuid"
	"time"
)

type CreateCustomerCommand struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateCustomerOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	err = uc.CustomerGateway.Create(customer)
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
