package postgres

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountPgGateway struct {
	DB *sql.DB
}

func NewAccountPgGateway(db *sql.DB) *AccountPgGateway {
	return &AccountPgGateway{DB: db}
}

func (a AccountPgGateway) Create(account *entity.Account) error {
	query := `INSERT INTO accounts (id, customer_id, balance, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5)`

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

func (a AccountPgGateway) UpdateBalance(ID uuid.UUID, amount decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2`

	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(amount, ID)
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
			  	WHERE id = $1`
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
