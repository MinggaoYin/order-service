package services

import "order-service/models"

type OrderService interface {
	GetById(id int64) (*models.Order, error)
	PlaceOrder(origins, destinations []string) (*models.Order, error)
	TakeOrder(order *models.Order) (*models.Order, error)
	ListOrders(offset, limit int) ([]models.Order, error)
}
