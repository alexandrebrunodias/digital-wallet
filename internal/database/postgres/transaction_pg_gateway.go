package postgres

import (
	"database/sql"
	"github.com/alexandrebrunodias/wallet-core/internal/entity"
	"github.com/google/uuid"
)

type TransactionPgGateway struct {
	DB *sql.DB
}

func NewTransactionPgGateway(db *sql.DB) *TransactionPgGateway {
	return &TransactionPgGateway{DB: db}
}

func (a TransactionPgGateway) Create(transaction *entity.Transaction) error {
	query := `INSERT INTO transactions (id, from_account_id, to_account_id, amount, created_at) 
				VALUES ($1, $2, $3, $4, $5)`
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		transaction.ID,
		transaction.FromAccount.ID,
		transaction.ToAccount.ID,
		transaction.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a TransactionPgGateway) GetByID(ID uuid.UUID) (*entity.Transaction, error) {
	var transaction entity.Transaction
	var fromAccount entity.Account
	var toAccount entity.Account
	transaction.FromAccount = &fromAccount
	transaction.ToAccount = &toAccount

	query := `SELECT id, from_account_id, to_account_id, amount, created_at
			  	FROM transactions
			  	WHERE id = $1`
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(ID).
		Scan(
			&transaction.ID,
			&transaction.FromAccount.ID,
			&transaction.ToAccount.ID,
			&transaction.Amount,
			&transaction.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &transaction, err
}
