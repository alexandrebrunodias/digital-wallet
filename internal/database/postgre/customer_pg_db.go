package postgre

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type CustomerPgDB struct {
	DB *sql.DB
}

func NewCustomerPgDB(db *sql.DB) *CustomerPgDB {
	return &CustomerPgDB{DB: db}
}

func (c *CustomerPgDB) Save(customer *entity.Customer) error {
	statement, err := c.DB.Prepare(
		"INSERT INTO customers (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}

	_, err = statement.Exec(customer.ID, customer.Name, customer.Email, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (c *CustomerPgDB) GetByID(ID uuid.UUID) (*entity.Customer, error) {
	customer := &entity.Customer{}

	statement, err := c.DB.Prepare(
		"SELECT id, name, email, created_at, updated_at FROM customers WHERE id = ?",
	)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	err = statement.
		QueryRow(ID).
		Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
