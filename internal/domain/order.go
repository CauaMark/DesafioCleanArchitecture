package domain

type Order struct {
	ID           int64   `json:"id"`
	CustomerName string  `json:"customerName"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
}

type OrderRepository interface {
	List() ([]Order, error)
	Create(order Order) (Order, error)
}
