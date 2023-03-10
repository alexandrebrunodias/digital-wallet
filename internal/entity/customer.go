package entity

import (
	"errors"
	"github.com/google/uuid"
	"net/mail"
	"strings"
	"time"
)

type Customer struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCustomer(name string, email string) (*Customer, error) {
	now := time.Now()
	customer := &Customer{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := customer.Validate()

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *Customer) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("'name' should not be blank")
	}

	if _, err := mail.ParseAddress(c.Email); err != nil {
		return errors.New("'email' is invalid")
	}

	return nil
}

func (c *Customer) Update(name string, email string) error {
	c.Name = name
	c.Email = email
	return c.Validate()
}
