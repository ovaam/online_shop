package models

type OrderStatus string

const (
	StatusNew       OrderStatus = "NEW"
	StatusFinished  OrderStatus = "FINISHED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
	Status      OrderStatus `json:"status"`
}
