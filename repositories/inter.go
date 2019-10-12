package repositories

import "order-service/models"

type OrderRepository interface {
	GetById(id int64) (*models.Order, error)
	Update(order *models.Order, withStatus string) (*models.Order, error)
	Create(o *models.Order) (*models.Order, error)
	Delete(id int64) (bool, error)
	List(offset, limit int) ([]models.Order, error)
}
