package expense

type UpdateDTO struct {
	Description *string          `json:"description,omitempty"`
	Amount      *float64         `json:"amount,omitempty"`
	Category    *Category        `json:"category,omitempty"`
	Type        *TransactionType `json:"transaction_type,omitempty"`
}
