package expense

type CreateDTO struct {
	Description string          `json:"description"`
	Amount      float64         `json:"amount"`
	Category    Category        `json:"categoty"`
	Type        TransactionType `json:"type"`
}
