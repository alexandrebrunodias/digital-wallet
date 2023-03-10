package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCustomer_CreateSuccessfully(t *testing.T) {
	expectedName := "Alex"
	expectedEmail := "alexandrebrunodias@gmail.com"
	customer, err := newCustomer(expectedName, expectedEmail)

	assert.Nil(t, err)
	assert.Equal(t, expectedName, customer.Name)
	assert.Equal(t, expectedEmail, expectedEmail)
}

func TestNewCustomer_ErrorDueToEmptyName(t *testing.T) {
	expectedName := ""
	expectedEmail := "alexandrebrunodias@gmail.com"
	customer, err := newCustomer(expectedName, expectedEmail)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "'name' should not be blank")
	assert.Nil(t, customer)
}

func TestNewCustomer_ErrorDueToEmptyEmail(t *testing.T) {
	expectedName := "Alex"
	expectedEmail := ""
	customer, err := newCustomer(expectedName, expectedEmail)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "'email' is invalid")
	assert.Nil(t, customer)
}

func TestNewCustomer_ErrorDueToInvalidEmail(t *testing.T) {
	expectedName := "Alex"
	expectedEmail := "invalid_email"
	customer, err := newCustomer(expectedName, expectedEmail)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "'email' is invalid")
	assert.Nil(t, customer)
}

func TestUpdateCustomer_UpdatedSuccessfully(t *testing.T) {
	expectedName := "Xela"
	expectedEmail := "erdnaxelaonurbsaid@giamg.moc"
	customer, _ := newCustomer("Alex", "alexandrebrunodias@gmail.com")

	err := customer.Update(expectedName, expectedEmail)

	assert.Nil(t, err)
	assert.Equal(t, expectedName, customer.Name)
	assert.Equal(t, expectedEmail, expectedEmail)
}

func TestUpdateCustomer_ErrorDueToInvalidParam(t *testing.T) {
	expectedName := ""
	expectedEmail := "erdnaxelaonurbsaid@giamg.moc"
	customer, _ := newCustomer("Alex", "alexandrebrunodias@gmail.com")

	err := customer.Update(expectedName, expectedEmail)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "'name' should not be blank")
}
