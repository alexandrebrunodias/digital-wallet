package postgre

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type AccountPgGateway struct {
	DB *sql.DB
}

func NewAccountPgGateway(db *sql.DB) *AccountPgGateway {
	return &AccountPgGateway{DB: db}
}

func (a AccountPgGateway) Save(account *entity.Account) error {
	query := `INSERT INTO accounts (id, customer_id, balance, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?)`

	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		account.ID,
		account.Customer.ID,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a AccountPgGateway) GetByID(ID uuid.UUID) (*entity.Account, error) {
	var account entity.Account
	var customer entity.Customer
	account.Customer = &customer

	query := `SELECT id, customer_id, balance, created_at, updated_at
			  	FROM accounts
			  	WHERE id = ?`
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(ID).
		Scan(
			&account.ID,
			&account.Customer.ID,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &account, err
}
