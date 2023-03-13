package postgre

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestNewCustomerPgDBTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerPgGatewaySuite))
}

func (s *CustomerPgGatewaySuite) TestSaveAndGetByID_SaveAndGetByIDSuccessfully() {
	expectedName := "alex"
	expectedEmail := "alexandrebrunodias@gmail.com"

	expectedCustomer, _ := entity.NewCustomer(expectedName, expectedEmail)

	err := s.CustomerPgGateway.Save(expectedCustomer)
	assert.Nil(s.T(), err)

	actualCustomer, err := s.CustomerPgGateway.GetByID(expectedCustomer.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), actualCustomer)
	assert.Equal(s.T(), expectedCustomer.ID, actualCustomer.ID)
	assert.Equal(s.T(), expectedCustomer.Name, actualCustomer.Name)
	assert.Equal(s.T(), expectedCustomer.Email, actualCustomer.Email)
	assert.Equal(s.T(), expectedCustomer.CreatedAt, actualCustomer.CreatedAt)
	assert.Equal(s.T(), expectedCustomer.UpdatedAt, actualCustomer.UpdatedAt)
}

func (s *CustomerPgGatewaySuite) TestGetByID_FetchEmpty() {
	actualCustomer, err := s.CustomerPgGateway.GetByID(uuid.New())
	expectedError := "sql: no rows in result set"

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedError, err.Error())
	assert.Nil(s.T(), actualCustomer)
}

type CustomerPgGatewaySuite struct {
	suite.Suite
	CustomerPgGateway *CustomerPgGatewayDB
}

func (s *CustomerPgGatewaySuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Require().Nil(err)

	s.CustomerPgGateway = NewCustomerPgGateway(db)

	stmt := `CREATE TABLE customers (
						id binary(16) PRIMARY KEY,
						name VARCHAR(255) NOT NULL,
						email VARCHAR(255) NOT NULL,
						created_at DATETIME,
						updated_at DATETIME
				   )`

	_, err = s.CustomerPgGateway.DB.Exec(stmt)
	s.Require().Nil(err)
}

func (s *CustomerPgGatewaySuite) TearDownSuite() {
	defer s.CustomerPgGateway.DB.Close()
	_, _ = s.CustomerPgGateway.DB.Exec("DROP TABLE customers")
}
