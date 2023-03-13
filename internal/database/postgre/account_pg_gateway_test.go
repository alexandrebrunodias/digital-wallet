package postgre

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestNewAccountPgDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountPgGatewaySuite))
}

func (s *AccountPgGatewaySuite) TestSaveAndGetByID_SaveSuccessfully() {
	err := s.AccountPgGateway.Save(s.AccountOne)
	assert.Nil(s.T(), err)

	actualAccount, err := s.AccountPgGateway.GetByID(s.AccountOne.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), actualAccount)
	assert.Equal(s.T(), s.AccountOne.ID, actualAccount.ID)
	assert.Equal(s.T(), s.AccountOne.Customer.ID, actualAccount.Customer.ID)
	assert.Equal(s.T(), s.AccountOne.Balance.String(), actualAccount.Balance.String())
	assert.Equal(s.T(), s.AccountOne.CreatedAt, actualAccount.CreatedAt)
	assert.Equal(s.T(), s.AccountOne.UpdatedAt, actualAccount.UpdatedAt)
}

func (s *AccountPgGatewaySuite) TestSave_FailDueInvalidAccount() {
	expectedPanicMessage := "runtime error: invalid memory address or nil pointer dereference"
	assert.Panicsf(s.T(), func() {
		_ = s.AccountPgGateway.Save(&entity.Account{})
	}, expectedPanicMessage)

	actualAccount, err := s.AccountPgGateway.GetByID(s.AccountOne.ID)

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), actualAccount)
}

func (s *AccountPgGatewaySuite) TestGetByID_GetSuccessfully() {
	_ = s.AccountPgGateway.Save(s.AccountOne)
	_ = s.AccountPgGateway.Save(s.AccountTwo)

	actualAccountOne, err := s.AccountPgGateway.GetByID(s.AccountOne.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), actualAccountOne)
	assert.Equal(s.T(), s.AccountOne.ID, actualAccountOne.ID)
	assert.Equal(s.T(), s.AccountOne.Customer.ID, actualAccountOne.Customer.ID)
	assert.Equal(s.T(), s.AccountOne.Balance.String(), actualAccountOne.Balance.String())
	assert.Equal(s.T(), s.AccountOne.CreatedAt, actualAccountOne.CreatedAt)
	assert.Equal(s.T(), s.AccountOne.UpdatedAt, actualAccountOne.UpdatedAt)

	actualAccountTwo, err := s.AccountPgGateway.GetByID(s.AccountTwo.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), actualAccountTwo)
	assert.Equal(s.T(), s.AccountTwo.ID, actualAccountTwo.ID)
	assert.Equal(s.T(), s.AccountTwo.Customer.ID, actualAccountTwo.Customer.ID)
	assert.Equal(s.T(), s.AccountTwo.Balance.String(), actualAccountTwo.Balance.String())
	assert.Equal(s.T(), s.AccountTwo.CreatedAt, actualAccountTwo.CreatedAt)
	assert.Equal(s.T(), s.AccountTwo.UpdatedAt, actualAccountTwo.UpdatedAt)
}

func (s *AccountPgGatewaySuite) TestGetByID_FetchEmpty() {
	actualAccount, err := s.AccountPgGateway.GetByID(uuid.New())
	expectedError := "sql: no rows in result set"

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedError, err.Error())
	assert.Nil(s.T(), actualAccount)
}

type AccountPgGatewaySuite struct {
	suite.Suite
	AccountPgGateway *AccountPgGateway
	AccountOne       *entity.Account
	AccountTwo       *entity.Account
	Customer         *entity.Customer
}

func (s *AccountPgGatewaySuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Require().Nil(err)

	s.AccountPgGateway = NewAccountPgGateway(db)

	query := `CREATE TABLE customers (
				id binary(16) PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				email VARCHAR(255) NOT NULL,
				created_at DATETIME,
				updated_at DATETIME
			 )`
	_, err = s.AccountPgGateway.DB.Exec(query)
	s.Require().Nil(err)

	query = `CREATE TABLE accounts (
				id BINARY(16) PRIMARY KEY,
				customer_id BINARY(16) NOT NULL,
				balance DECIMAL(12, 2),
				created_at DATETIME,
				updated_at DATETIME
		     )`

	_, err = s.AccountPgGateway.DB.Exec(query)
	s.Require().Nil(err)

	s.Customer, err = entity.NewCustomer("alex", "alexandrebrunodias@gmail.com")
	s.Require().Nil(err)

	err = NewCustomerPgGateway(db).Save(s.Customer)
	s.Require().Nil(err)
}

func (s *AccountPgGatewaySuite) SetupTest() {
	var err error
	s.AccountOne, err = entity.NewAccount(s.Customer)
	s.Require().Nil(err)

	s.AccountTwo, err = entity.NewAccount(s.Customer)
	s.Require().Nil(err)
}

func (s *AccountPgGatewaySuite) TearDownSuite() {
	defer s.AccountPgGateway.DB.Close()
	_, _ = s.AccountPgGateway.DB.Exec("DROP TABLE accounts")
	_, _ = s.AccountPgGateway.DB.Exec("DROP TABLE customers")
}
