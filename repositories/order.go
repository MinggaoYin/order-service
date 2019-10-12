package repositories

import (
	"database/sql"

	"order-service/models"
)

type OrderRepo struct {
	Conn *sql.DB
}

func NewMysqlOrderRepo(conn *sql.DB) *OrderRepo {
	return &OrderRepo{conn}
}

func (rp *OrderRepo) fetch(query string, args ...interface{}) ([]models.Order, error) {
	rows, err := rp.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}

	orders := make([]models.Order, 0)

	for rows.Next() {
		order := models.Order{}

		err := rows.Scan(&order.Id, &order.Distance, &order.Status)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (rp *OrderRepo) GetById(id int64) (*models.Order, error) {
	query := "SELECT id, distance, status FROM orders WHERE id = ?"
	orders, err := rp.fetch(query, id)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, models.ErrNotFound
	}

	return &orders[0], nil
}

func (rp *OrderRepo) Update(order *models.Order, withStatus string) (*models.Order, error) {
	query := "UPDATE orders SET status = ? where id = ? AND status = ?"

	stmt, err := rp.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(
		order.Status,
		order.Id,
		withStatus)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected != 1 {
		return nil, models.ErrCannotUpdate
	}

	return order, nil
}

func (rp *OrderRepo) Create(order *models.Order) (*models.Order, error) {
	query := "INSERT INTO orders (origin_lat, origin_lng, destination_lat, destination_lng, distance, status) VALUES (?, ?, ?, ?, ?, ?)"

	stmt, err := rp.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(
		order.Origins[0],
		order.Origins[1],
		order.Destinations[0],
		order.Destinations[1],
		order.Distance,
		order.Status)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	order.Id = id
	return order, nil
}

func (rp *OrderRepo) Delete(id int64) (bool, error) {
	query := "DELETE FROM orders WHERE id = ?"

	stmt, err := rp.Conn.Prepare(query)
	if err != nil {
		return false, err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rp *OrderRepo) List(offset, limit int) ([]models.Order, error) {
	query := "SELECT id, distance, status FROM orders ORDER BY id ASC LIMIT ?, ?"

	orders, err := rp.fetch(query, offset, limit)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
