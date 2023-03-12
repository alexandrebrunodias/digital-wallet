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

func TestNewCustomerPGDBTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerPgDBSuite))
}

func (s *CustomerPgDBSuite) TestSaveAndGetByID_SaveAndGetByIDSuccessfully() {
	expectedName := "alex"
	expectedEmail := "alexandrebrunodias@gmail.com"

	expectedCustomer, _ := entity.NewCustomer(expectedName, expectedEmail)

	err := s.CustomerPgDB.Save(expectedCustomer)
	assert.Nil(s.T(), err)

	actualCustomer, err := s.CustomerPgDB.GetByID(expectedCustomer.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), actualCustomer)
	assert.Equal(s.T(), expectedCustomer.ID, actualCustomer.ID)
	assert.Equal(s.T(), expectedCustomer.Name, actualCustomer.Name)
	assert.Equal(s.T(), expectedCustomer.Email, actualCustomer.Email)
	assert.Equal(s.T(), expectedCustomer.CreatedAt, actualCustomer.CreatedAt)
	assert.Equal(s.T(), expectedCustomer.UpdatedAt, actualCustomer.UpdatedAt)
}

func (s *CustomerPgDBSuite) TestGetByID_FetchEmpty() {
	actualCustomer, err := s.CustomerPgDB.GetByID(uuid.New())
	expectedError := "sql: no rows in result set"

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), expectedError, err.Error())
	assert.Nil(s.T(), actualCustomer)
}

type CustomerPgDBSuite struct {
	suite.Suite
	CustomerPgDB *CustomerPgDB
}

func (s *CustomerPgDBSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Require().Nil(err)

	s.CustomerPgDB = NewCustomerPgDB(db)

	statement := `
					CREATE TABLE customers (
					    id binary(16) PRIMARY KEY,
					    name VARCHAR(255) NOT NULL,
					    email VARCHAR(255) NOT NULL,
					    created_at DATETIME,
					    updated_at DATETIME
					)
				 `
	_, err = s.CustomerPgDB.DB.Exec(statement)
	s.Require().Nil(err)

}

func (s *CustomerPgDBSuite) TearDownSuite() {
	defer s.CustomerPgDB.DB.Close()
	s.CustomerPgDB.DB.Exec("DROP TABLE customers")
}
