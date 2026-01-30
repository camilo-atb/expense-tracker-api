package expense

import "time"

// Transaction represents a financial movement, either income or expense.
type Transaction struct {
	ID          int
	Description string
	Amount      float64
	Category    Category
	Type        TransactionType
	Date        time.Time
	Status      Status
}

type Category string

const (
	Comida          Category = "comida"
	Transporte      Category = "transporte"
	Ocio            Category = "ocio"
	Entretenimiento Category = "entretenimiento"
	Estudios        Category = "estudios"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Status string

const (
	Active  Status = "active"
	Removed Status = "removed"
)

// --

func (t TransactionType) IsValid() bool {
	return t == Income || t == Expense
}

func (c Category) IsValid() bool {
	switch c {
	case Comida, Transporte, Ocio, Entretenimiento, Estudios:
		return true
	default:
		return false
	}
}
