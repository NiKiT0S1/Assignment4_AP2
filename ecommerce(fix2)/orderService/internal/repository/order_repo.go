package repository

import (
	"orderService/internal/domain"

	"github.com/jmoiron/sqlx"
)

type orderRepo struct {
	db *sqlx.DB
}

// 4.3) Начинаем транзакцию для атомарного создания заказа(либо операция полностью успешно завершится, либо нет) со всеми его элементами
func NewOrderRepo(db *sqlx.DB) domain.OrderRepository {
	return &orderRepo{db}
}

func (r *orderRepo) Create(order *domain.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// 4.4) Создаем запись в таблице orders и получаем ID
	var orderID int
	err = tx.QueryRowx(`
		INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id
	`, order.UserID, order.Status).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`, orderID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 4.5) Фиксируем транзакцию и устанавливаем полученный ID в объект заказа
	tx.Commit()
	order.ID = orderID
	return nil
}

func (r *orderRepo) GetByID(id int) (*domain.Order, error) {
	var o domain.Order
	err := r.db.Get(&o, "SELECT * FROM orders WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	err = r.db.Select(&o.Items, "SELECT * FROM order_items WHERE order_id=$1", o.ID)
	return &o, err
}

func (r *orderRepo) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec("UPDATE orders SET status=$1 WHERE id=$2", status, id)
	return err
}

func (r *orderRepo) ListByUser(userID int) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Select(&orders, "SELECT * FROM orders WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}

	for i, order := range orders {
		r.db.Select(&orders[i].Items, "SELECT * FROM order_items WHERE order_id=$1", order.ID)
	}

	return orders, nil
}
