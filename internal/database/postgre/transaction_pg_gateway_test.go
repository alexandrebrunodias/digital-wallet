package postgre

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestNewTransactionPgDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionPgGatewaySuite))
}

func (s *TransactionPgGatewaySuite) TestSave_SaveSuccessfully() {
	fromAccountInitialBalance := decimal.NewFromInt(2000)
	expectedAmount := decimal.NewFromInt(1000)

	_ = s.FromAccount.Credit(fromAccountInitialBalance)

	expectedTransaction, _ := entity.NewTransaction(s.FromAccount, s.ToAccount, expectedAmount)
	err := s.TransactionPgGateway.Save(expectedTransaction)
	assert.Nil(s.T(), err)

	actualTransaction, err := s.TransactionPgGateway.GetByID(expectedTransaction.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), actualTransaction)
	assert.Equal(s.T(), expectedTransaction.ID, actualTransaction.ID)
	assert.Equal(s.T(), expectedAmount, actualTransaction.Amount)
	assert.Equal(s.T(), s.FromAccount.ID, actualTransaction.FromAccount.ID)
	assert.Equal(s.T(), s.ToAccount.ID, actualTransaction.ToAccount.ID)
	assert.Equal(s.T(), expectedTransaction.CreatedAt, actualTransaction.CreatedAt)
}

func (s *TransactionPgGatewaySuite) TestSave_FailDueInvalidAccount() {
	expectedPanicMessage := "runtime error: invalid memory address or nil pointer dereference"
	assert.Panicsf(s.T(), func() {
		_ = s.TransactionPgGateway.Save(&entity.Transaction{})
	}, expectedPanicMessage)

	actualAccount, err := s.TransactionPgGateway.GetByID(s.FromAccount.ID)

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), actualAccount)
}

func (s *TransactionPgGatewaySuite) TestGetByID_FetchEmpty() {
	actualAccount, err := s.TransactionPgGateway.GetByID(uuid.New())
	expectedError := "sql: no rows in result set"

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedError, err.Error())
	assert.Nil(s.T(), actualAccount)
}

type TransactionPgGatewaySuite struct {
	suite.Suite
	TransactionPgGateway *TransactionPgGateway
	FromAccount          *entity.Account
	ToAccount            *entity.Account
}

func (s *TransactionPgGatewaySuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Require().Nil(err)

	s.TransactionPgGateway = NewTransactionPgGateway(db)

	query := `CREATE TABLE transactions (
				id BINARY(16) PRIMARY KEY,
				from_account_id BINARY(16) NOT NULL,
				to_account_id BINARY(16) NOT NULL,
				amount DECIMAL(14, 2),
				created_at DATETIME
		     )`

	_, err = s.TransactionPgGateway.DB.Exec(query)
	s.Require().Nil(err)
}

func (s *TransactionPgGatewaySuite) SetupTest() {
	customer, err := entity.NewCustomer("alex", "alexandrebrunodias@gmail.com")
	s.Require().Nil(err)

	s.FromAccount, err = entity.NewAccount(customer)
	s.Require().Nil(err)

	s.ToAccount, err = entity.NewAccount(customer)
	s.Require().Nil(err)
}

func (s *TransactionPgGatewaySuite) TearDownSuite() {
	defer s.TransactionPgGateway.DB.Close()
	_, _ = s.TransactionPgGateway.DB.Exec("DROP TABLE transactions")
}
