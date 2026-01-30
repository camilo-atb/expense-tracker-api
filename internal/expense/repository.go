package expense

import (
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	AddTransaction(createDTO *CreateDTO) (*Transaction, error)
	ModifyTransaction(id int, transaccionDto *UpdateDTO) (*Transaction, error)
	DeleteTransaction(id int) error
	DeleteLogicalTransaction(id int) error
	GetAllTransactions() ([]*Transaction, error)
	GetTransactionById(id int) (*Transaction, error)
	GetTransactionsByType(typeT TransactionType) ([]*Transaction, error)
	GetTransactionsByTypeAndDate(typeT TransactionType, from time.Time, to time.Time) ([]*Transaction, error)
	GetTransactionsBetweenDates(from time.Time, to time.Time) ([]*Transaction, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(dbIny *sql.DB) Repository {
	return &repository{db: dbIny}
}

var ErrRepoNotFound = errors.New("repo: not found")

const transactionColumns = "id, description, amount, category, type, date, status"

func (r *repository) AddTransaction(createDTO *CreateDTO) (*Transaction, error) {
	q := "INSERT INTO transactions (description, amount, category, type) VALUES ($1, $2, $3, $4) RETURNING " + transactionColumns
	result := r.db.QueryRow(q, createDTO.Description, createDTO.Amount, createDTO.Category, createDTO.Type)
	transaction, err := r.scanTransaction(result)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *repository) ModifyTransaction(id int, updateDTO *UpdateDTO) (*Transaction, error) {
	q := "UPDATE transactions SET description = COALESCE($1, description), amount = COALESCE($2, amount), category = COALESCE($3, category), type = COALESCE($4, type) WHERE id = $5 RETURNING " + transactionColumns
	row := r.db.QueryRow(q, updateDTO.Description, updateDTO.Amount, updateDTO.Category, updateDTO.Type, id)
	transaction, err := r.scanTransaction(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRepoNotFound
		}
		return nil, err
	}
	return transaction, nil
}

func (r *repository) DeleteTransaction(id int) error {
	q := "DELETE FROM transactions WHERE id = $1"

	result, err := r.db.Exec(q, id)

	if err != nil {
		return err
	}

	response, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if response == 0 {
		return ErrRepoNotFound
	}

	return nil
}

func (r *repository) DeleteLogicalTransaction(id int) error {
	q := "UPDATE transactions SET status = $1 WHERE id = $2"

	result, err := r.db.Exec(q, Removed, id)

	if err != nil {
		return err
	}

	response, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if response == 0 {
		return ErrRepoNotFound
	}

	return nil
}

func (r *repository) GetAllTransactions() ([]*Transaction, error) {
	q := "SELECT " + transactionColumns + " FROM transactions"
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*Transaction{}

	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *repository) GetTransactionById(id int) (*Transaction, error) {
	q := "SELECT " + transactionColumns + " FROM transactions WHERE id = $1"
	row := r.db.QueryRow(q, id)
	transaction, err := r.scanTransaction(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRepoNotFound
		}
		return nil, err
	}

	return transaction, nil
}

func (r *repository) GetTransactionsByType(typeT TransactionType) ([]*Transaction, error) {
	q := "SELECT " + transactionColumns + " FROM transactions WHERE type = $1"
	rows, err := r.db.Query(q, typeT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactions := []*Transaction{}
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *repository) GetTransactionsByTypeAndDate(typeT TransactionType, from time.Time, to time.Time) ([]*Transaction, error) {
	q := "SELECT " + transactionColumns + " FROM transactions WHERE type = $1 AND date BETWEEN $2 AND $3"
	rows, err := r.db.Query(q, typeT, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactions := []*Transaction{}
	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *repository) GetTransactionsBetweenDates(from time.Time, to time.Time) ([]*Transaction, error) {
	q := "SELECT " + transactionColumns + " FROM transactions WHERE date BETWEEN $1 AND $2"
	rows, err := r.db.Query(q, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions = []*Transaction{}

	for rows.Next() {
		transaction, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

//--

func (r *repository) scanTransaction(scanner interface{ Scan(dest ...any) error }) (*Transaction, error) {
	var t Transaction
	err := scanner.Scan(&t.ID, &t.Description, &t.Amount, &t.Category, &t.Type, &t.Date, &t.Status)
	return &t, err
}
