package services

import (
	"strconv"
	"strings"

	mdorder "order-service/models/order"

	"github.com/sirupsen/logrus"
)

func PlaceOrder(origin, destination []string) (*mdorder.Order, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service", "method": "PlaceOrder"})

	distance, err := getDistance(
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

	order, err := mdorder.Create(origins, destinations, distance, mdorder.StatusUnassigned)
	if err != nil {
		log.WithError(err).Error("Failed to create order")
		return nil, err
	}

	return order, nil
}

func TakeOrder(id int64) error {
	log := logrus.WithFields(logrus.Fields{"module": "service", "method": "TakeOrder", "order_id": id})

	result, err := mdorder.Exec("UPDATE orders SET status = ? WHERE id = ? AND status = ?",
		mdorder.StatusTaken, id, mdorder.StatusUnassigned)
	if err != nil {
		log.WithError(err).Error("Failed to take order")
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Failed to take order")
		return err
	}

	if rowCount == 0 {
		return ErrOrderAlreadyTaken
	}

	return nil
}

func ListOrders(offset, limit int) ([]mdorder.Order, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service", "method": "TakeOrder", "offset": offset, "limit": limit})

	orders, err := mdorder.List(offset, limit)
	if err != nil {
		log.WithError(err).Error("Failed to list orders")
		return nil, err
	}

	return orders, nil
}
