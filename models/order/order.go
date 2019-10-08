package order

import (
	"database/sql"

	"order-service/models"
)

var (
	StatusUnassigned = "UNASSIGNED"
	StatusTaken      = "TAKEN"
)

type Order struct {
	Id           int64    `json:"id"`
	Origins      []string `json:"-"`
	Destinations []string `json:"-"`
	Distance     int      `json:"distance"`
	Status       string   `json:"status"`
}

var createTableStat = `CREATE TABLE IF NOT EXISTS orders (
    id BIGINT(20) UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
    origin_lat DOUBLE NOT NULL,
    origin_lng DOUBLE NOT NULL,
    destination_lat DOUBLE NOT NULL,
    destination_lng DOUBLE NOT NULL,
    status VARCHAR(20) NOT NULL,
    distance INT UNSIGNED NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;
`

func init() {
	// create order table if not exists
	models.Db.Exec(createTableStat)
}

func Create(origins, destinations []float64, distance int, status string) (*Order, error) {
	result, err := models.Db.Exec("INSERT INTO orders (origin_lat, origin_lng, destination_lat, destination_lng, distance, status) VALUES (?, ?, ?, ?, ?, ?)",
		origins[0], origins[1], destinations[0], destinations[1], distance, status)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	order := &Order{
		Id:       id,
		Distance: distance,
		Status:   status,
	}

	return order, nil
}

func Get(id int64) (*Order, error) {
	var o Order
	err := models.Db.QueryRow("SELECT id, distance, status FROM orders where id = ?", id).
		Scan(&o.Id, &o.Distance, &o.Status)
	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func Update(o Order) error {
	_, err := models.Db.Exec("UPDATE orders SET status = ? where id = ?", o.Status, o.Id)
	return err
}

func Delete(id string) error {
	var o Order
	err := models.Db.QueryRow("DELETE FROM orders WHERE id = ?", id).
		Scan(&o.Id, &o.Distance, o.Status)

	return err
}

func List(offset, limit int) ([]Order, error) {
	results, err := models.Db.Query("SELECT id, distance, status FROM orders ORDER BY id ASC LIMIT ?, ?", offset, limit)
	if err != nil {
		return nil, err
	}

	orders := []Order{}

	for results.Next() {
		var o Order
		err = results.Scan(&o.Id, &o.Distance, &o.Status)
		if err != nil {
			continue
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := models.Db.Exec(query, args...)
	return result, err
}
