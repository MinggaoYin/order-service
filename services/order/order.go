package order

import (
	"strconv"
	"strings"

	"order-service/models"
	"order-service/repositories"
	"order-service/services"

	"github.com/sirupsen/logrus"
)

type orderService struct {
	orderRepo          repositories.OrderRepository
	distanceCalculator services.DistanceCalculator
}

// NewOrderService will create new an OrderService object representation of OrderService interface
func NewOrderService(o repositories.OrderRepository, distanceCalculator services.DistanceCalculator) services.OrderService {
	return &orderService{
		orderRepo:          o,
		distanceCalculator: distanceCalculator,
	}
}

func (s *orderService) GetById(id int64) (*models.Order, error) {
	return s.orderRepo.GetById(id)
}

func (s *orderService) PlaceOrder(origin, destination []string) (*models.Order, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service/order", "method": "PlaceOrder"})

	distance, err := s.distanceCalculator.GetDistance(
		[]string{strings.Join(origin, ",")},
		[]string{strings.Join(destination, ",")})
	if err != nil {
		log.WithError(err).Error("Failed to get distance")
		return nil, err
	}

	origins := []float64{}
	originLat, _ := strconv.ParseFloat(origin[0], 64)
	originLog, _ := strconv.ParseFloat(origin[1], 64)
	origins = append(origins, originLat)
	origins = append(origins, originLog)

	destinations := []float64{}
	destLat, _ := strconv.ParseFloat(destination[0], 64)
	destLog, _ := strconv.ParseFloat(destination[1], 64)
	destinations = append(destinations, destLat)
	destinations = append(destinations, destLog)

	order := &models.Order{
		Origins:      origins,
		Destinations: destinations,
		Distance:     distance,
		Status:       models.StatusUnassigned,
	}

	order, err = s.orderRepo.Create(order)
	if err != nil {
		log.WithError(err).Error("Failed to create order")
		return nil, err
	}

	return order, nil
}

func (s *orderService) TakeOrder(order *models.Order) (*models.Order, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service/order", "method": "TakeOrder", "order_id": order.Id, "current_status": order.Status})

	if order.Status != models.StatusUnassigned {
		log.Debug("Failed to take order, since it was already taken")
		return nil, services.ErrOrderAlreadyTaken
	}

	order.Status = models.StatusTaken
	_, err := s.orderRepo.Update(order, models.StatusUnassigned)
	if err == models.ErrCannotUpdate {
		log.Debug("Failed to take order, since it was already taken")
		return nil, services.ErrOrderAlreadyTaken
	} else if err != nil {
		log.WithError(err).Error("Failed to take order")
		return nil, err
	}

	return order, nil
}

func (s *orderService) ListOrders(offset, limit int) ([]models.Order, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service/order", "method": "ListOrders", "offset": offset, "limit": limit})

	orders, err := s.orderRepo.List(offset, limit)
	if err != nil {
		log.WithError(err).Error("Failed to list orders")
		return nil, err
	}

	return orders, nil
}
