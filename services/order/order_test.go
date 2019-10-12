package order

import (
	"errors"
	"order-service/services"
	"strings"
	"testing"

	"order-service/models"
	rpmocks "order-service/repositories/mocks"
	srvmocks "order-service/services/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_GetById(t *testing.T) {
	mockOrderRepo := new(rpmocks.OrderRepository)
	mockDistanceSrv := new(srvmocks.DistanceCalculator)

	mockOrder := &models.Order{
		Id:           1,
		Origins:      nil,
		Destinations: nil,
		Distance:     10,
		Status:       models.StatusUnassigned,
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("GetById", mock.AnythingOfType("int64")).Return(mockOrder, nil).Once()
		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		order, err := orderService.GetById(mockOrder.Id)
		assert.NoError(t, err)
		assert.NotNil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockOrderRepo.On("GetById", mock.AnythingOfType("int64")).Return(nil, errors.New("exception")).Once()
		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		order, err := orderService.GetById(mockOrder.Id)
		assert.Error(t, err)
		assert.Nil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_PlaceOrder(t *testing.T) {
	mockOrderRepo := new(rpmocks.OrderRepository)
	mockDistanceSrv := new(srvmocks.DistanceCalculator)

	origins := []float64{22.286681, 114.193260}
	destinations := []float64{22.279707, 114.186301}

	originStrs := []string{"22.286681", "114.193260"}
	destinationStrs := []string{"22.279707", "114.186301"}

	mockOrder := &models.Order{
		Id:           0,
		Origins:      origins,
		Destinations: destinations,
		Distance:     100,
		Status:       "UNASSIGNED",
	}

	t.Run("success", func(t *testing.T) {
		mockDistanceSrv.On("GetDistance", []string{strings.Join(originStrs, ",")}, []string{strings.Join(destinationStrs, ",")}).
			Return(100, nil).Once()
		mockOrderRepo.On("Create", mockOrder).Return(mockOrder, nil).Once()

		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		order, err := orderService.PlaceOrder(originStrs, destinationStrs)
		assert.NoError(t, err)
		assert.NotNil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("on get distance failed", func(t *testing.T) {
		mockDistanceSrv.On("GetDistance", []string{strings.Join(originStrs, ",")}, []string{strings.Join(destinationStrs, ",")}).
			Return(0, errors.New("exception")).Once()

		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		order, err := orderService.PlaceOrder(originStrs, destinationStrs)
		assert.Error(t, err)
		assert.Nil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("cannot get distance too far away", func(t *testing.T) {
		mockDistanceSrv.On("GetDistance", []string{strings.Join(originStrs, ",")}, []string{strings.Join(destinationStrs, ",")}).
			Return(0, services.ErrCannotCalculateDistance).Once()

		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		order, err := orderService.PlaceOrder(originStrs, destinationStrs)
		assert.Error(t, err)
		assert.EqualError(t, err, services.ErrCannotCalculateDistance.Error())
		assert.Nil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("cannot create order", func(t *testing.T) {
		mockDistanceSrv.On("GetDistance", []string{strings.Join(originStrs, ",")}, []string{strings.Join(destinationStrs, ",")}).
			Return(100, nil).Once()
		mockOrderRepo.On("Create", mockOrder).Return(nil, errors.New("exception")).Once()

		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		order, err := orderService.PlaceOrder(originStrs, destinationStrs)
		assert.Error(t, err)
		assert.Nil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_TakeOrder(t *testing.T) {
	mockOrderRepo := new(rpmocks.OrderRepository)
	mockDistanceSrv := new(srvmocks.DistanceCalculator)

	mockOrder := &models.Order{
		Id:           1,
		Origins:      nil,
		Destinations: nil,
		Distance:     10,
		Status:       models.StatusUnassigned,
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("Update", mockOrder, models.StatusUnassigned).Return(mockOrder, nil).Once()

		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		mockOrder.Status = models.StatusUnassigned
		order, err := orderService.TakeOrder(mockOrder)
		assert.NoError(t, err)
		assert.NotNil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order already taken before querying db", func(t *testing.T) {
		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		mockOrder.Status = models.StatusTaken
		order, err := orderService.TakeOrder(mockOrder)
		assert.Error(t, err)
		assert.EqualError(t, err, services.ErrOrderAlreadyTaken.Error())
		assert.Nil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order already taken after querying db", func(t *testing.T) {
		mockOrderRepo.On("Update", mockOrder, models.StatusUnassigned).Return(nil, models.ErrCannotUpdate).Once()

		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		mockOrder.Status = models.StatusUnassigned
		order, err := orderService.TakeOrder(mockOrder)
		assert.Error(t, err)
		assert.EqualError(t, err, services.ErrOrderAlreadyTaken.Error())
		assert.Nil(t, order)

		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_ListOrders(t *testing.T) {
	mockOrderRepo := new(rpmocks.OrderRepository)
	mockDistanceSrv := new(srvmocks.DistanceCalculator)

	origins := []float64{22.286681, 114.193260}
	destinations := []float64{22.279707, 114.186301}

	mockOrders := []models.Order{
		models.Order{
			Id:           1,
			Origins:      origins,
			Destinations: destinations,
			Distance:     10,
			Status:       models.StatusUnassigned,
		},
		models.Order{
			Id:           2,
			Origins:      origins,
			Destinations: destinations,
			Distance:     20,
			Status:       models.StatusTaken,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("List", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockOrders, nil).Once()
		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		orders, err := orderService.ListOrders(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, orders)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockOrderRepo.On("List", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errors.New("exception")).Once()
		orderService := NewOrderService(mockOrderRepo, mockDistanceSrv)

		orders, err := orderService.ListOrders(1, 1)
		assert.Error(t, err)
		assert.Nil(t, orders)

		mockOrderRepo.AssertExpectations(t)
	})
}
