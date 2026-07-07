package postgres

import (
	"database/sql"
	"desafio-clean-architecture/internal/domain"
	"fmt"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) List() ([]domain.Order, error) {
	rows, err := r.db.Query(`
		SELECT id, customer_name, total, status
		FROM orders
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.Total, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *OrderRepository) Create(order domain.Order) (domain.Order, error) {
	stmt := `
		INSERT INTO orders (customer_name, total, status)
		VALUES ($1, $2, $3)
		RETURNING id, customer_name, total, status
	`
	var created domain.Order
	err := r.db.QueryRow(stmt, order.CustomerName, order.Total, order.Status).Scan(
		&created.ID,
		&created.CustomerName,
		&created.Total,
		&created.Status,
	)
	if err != nil {
		return domain.Order{}, fmt.Errorf("create order: %w", err)
	}
	return created, nil
}
