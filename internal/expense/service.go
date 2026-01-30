package expense

import (
	"errors"
	"strings"
	"time"
)

type Service interface {
	AddTransaction(createDTO *CreateDTO) (*Transaction, error)
	ModifyTransaction(id int, transaccionDto *UpdateDTO) (*Transaction, error)
	DeleteTransaction(id int) error
	DeleteLogicalTransaction(id int) error
	GetAllTransactions() ([]*Transaction, error)
	GetTransactionById(id int) (*Transaction, error)
	GetTransactionsByType(typeT TransactionType) ([]*Transaction, error)
	GetTransactionsByTypeAndDate(typeT TransactionType, from string, to string) ([]*Transaction, error)
	GetTransactionsBetweenDates(from string, to string) ([]*Transaction, error)
	ObtainTotalByTypeAndDate(t TransactionType, from string, to string) (*float64, error)
	GetNetIncome() (*float64, error)
}

type service struct {
	repo Repository
}

func NewService(repoIny Repository) Service {
	return &service{repo: repoIny}
}

var (
	errInvalidInput        = errors.New("service: invalid input")
	errInvalidId           = errors.New("service: invalid id")
	ErrTransactionNotFound = errors.New("service: transaction not found")
)

func (s *service) AddTransaction(createDTO *CreateDTO) (*Transaction, error) {
	if strings.TrimSpace(createDTO.Description) == "" || createDTO.Amount == 0 {
		return nil, errInvalidInput
	}

	if createDTO.Amount <= 0 {
		return nil, errInvalidInput
	}

	if !createDTO.Category.IsValid() {
		return nil, errInvalidInput
	}

	if !createDTO.Type.IsValid() {
		return nil, errInvalidInput
	}

	transaction, err := s.repo.AddTransaction(createDTO)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *service) ModifyTransaction(id int, transaccionDto *UpdateDTO) (*Transaction, error) {
	if id <= 0 {
		return nil, errInvalidId
	}

	if transaccionDto.Description != nil {
		if strings.TrimSpace(*transaccionDto.Description) == "" {
			return nil, errInvalidInput
		}
	}

	if transaccionDto.Amount != nil {
		if *transaccionDto.Amount == 0 || *transaccionDto.Amount <= 0 {
			return nil, errInvalidInput
		}
	}

	if transaccionDto.Type != nil && !transaccionDto.Type.IsValid() {
		return nil, errInvalidInput
	}

	if transaccionDto.Category != nil && !transaccionDto.Category.IsValid() {
		return nil, errInvalidInput
	}

	transaction, err := s.repo.ModifyTransaction(id, transaccionDto)
	if err != nil {
		if errors.Is(err, ErrRepoNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	return transaction, nil
}

func (s *service) DeleteTransaction(id int) error {
	if id <= 0 {
		return errInvalidId
	}

	err := s.repo.DeleteTransaction(id)
	if err != nil {
		if errors.Is(err, ErrRepoNotFound) {
			return ErrTransactionNotFound
		}
		return err
	}
	return nil
}

func (s *service) DeleteLogicalTransaction(id int) error {
	if id <= 0 {
		return errInvalidId
	}

	err := s.repo.DeleteLogicalTransaction(id)
	if err != nil {
		if errors.Is(err, ErrRepoNotFound) {
			return ErrTransactionNotFound
		}
		return err
	}
	return nil
}

func (s *service) GetAllTransactions() ([]*Transaction, error) {
	transactions, err := s.repo.GetAllTransactions()
	if err != nil {
		return nil, err
	}
	return transactions, err
}

func (s *service) GetTransactionById(id int) (*Transaction, error) {
	if id <= 0 {
		return nil, errInvalidId
	}
	transaction, err := s.repo.GetTransactionById(id)

	if err != nil {
		if errors.Is(err, ErrRepoNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return transaction, nil
}

func (s *service) GetTransactionsByType(typeT TransactionType) ([]*Transaction, error) {
	if !typeT.IsValid() {
		return nil, errInvalidInput
	}

	transactions, err := s.repo.GetTransactionsByType(typeT)

	if err != nil {
		return nil, err
	}

	return transactions, err
}

func (s *service) GetTransactionsByTypeAndDate(typeT TransactionType, from string, to string) ([]*Transaction, error) {
	if !typeT.IsValid() {
		return nil, errInvalidInput
	}

	dateFrom, dateTo, err := s.parseDateRange(from, to)

	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.GetTransactionsByTypeAndDate(typeT, dateFrom, dateTo)

	if err != nil {
		return nil, err
	}

	return transactions, err
}

func (s *service) GetTransactionsBetweenDates(from string, to string) ([]*Transaction, error) {
	dateFrom, dateTo, err := s.parseDateRange(from, to)

	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.GetTransactionsBetweenDates(dateFrom, dateTo)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *service) ObtainTotalByTypeAndDate(t TransactionType, from string, to string) (*float64, error) {
	if !t.IsValid() {
		return nil, errInvalidInput
	}

	dateFrom, dateTo, err := s.parseDateRange(from, to)

	if err != nil {
		return nil, err
	}

	tx, err := s.repo.GetTransactionsByTypeAndDate(t, dateFrom, dateTo)

	if err != nil {
		return nil, err
	}

	total := s.calculateTotal(tx)

	return &total, nil
}

func (s *service) GetNetIncome() (*float64, error) {
	transactionsExpense, err := s.repo.GetTransactionsByType(Expense)

	if err != nil {
		return nil, err
	}

	transactionsIncome, err := s.repo.GetTransactionsByType(Income)

	if err != nil {
		return nil, err
	}

	netIncome := s.calculateTotal(transactionsIncome) - s.calculateTotal(transactionsExpense)

	return &netIncome, nil
}

// --------

func convertDate(date string) (dateTime time.Time, err error) {
	layout := "2006-01-02"
	parseDate, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, errInvalidInput
	}
	return parseDate, nil
}

func (s *service) parseDateRange(from, to string) (time.Time, time.Time, error) {
	df, err := convertDate(from)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	dt, err := convertDate(to)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return df, dt, nil
}

func (s *service) calculateTotal(transactions []*Transaction) float64 {
	var total float64
	for _, v := range transactions {
		total += v.Amount
	}
	return total
}
