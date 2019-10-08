package handlers

import (
	mdorder "order-service/models/order"
	srvorder "order-service/services"

	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// add custom validation for latitude and longitude
	validate.RegisterValidation("location", func(fl validator.FieldLevel) bool {
		latitudeStr := fl.Field().Index(0).String()
		longitudeStr := fl.Field().Index(1).String()

		return validate.Var(latitudeStr, "latitude") == nil &&
			validate.Var(longitudeStr, "longitude") == nil
	})
}

type PlaceOrderReq struct {
	Origin      []string `json:"origin" validate:"required,len=2,location"`
	Destination []string `json:"destination" validate:"required,len=2,location"`
}

func PlaceOrder(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "PlaceOrder"})

	var req PlaceOrderReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	log = log.WithFields(logrus.Fields{"origin": req.Origin, "destination": req.Destination})

	err = validate.Struct(req)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	order, err := srvorder.PlaceOrder(req.Origin, req.Destination)
	if err == srvorder.ErrCannotCalculateDistance {
		log.WithField("err", err).Error("Failed to calculate distance for location")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": srvorder.ErrCannotCalculateDistance.Error(),
		})
		return
	}
	if err != nil {
		log.WithField("err", err).Error("Failed to place order")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	log.WithField("order_id", order.Id).Debug("Successfully created order")
	ctx.JSON(order)
}

type TakeOrderReq struct {
	Status string `json:"status" validate:"required,eq=TAKEN"`
}

func TakeOrder(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "TakeOrder"})

	var req TakeOrderReq
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": "Invalid json provided",
		})
		return
	}

	err = validate.Struct(req)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	order := ctx.Values().Get("_order").(*mdorder.Order)

	log = log.WithFields(logrus.Fields{"order_id": order.Id, "status": order.Status})

	if order.Status != mdorder.StatusUnassigned {
		log.Debug("Failed to take order, since order already taken")
		ctx.StatusCode(iris.StatusConflict)
		ctx.JSON(iris.Map{
			"error": "Order already taken",
		})
		return
	}

	err = srvorder.TakeOrder(order.Id)
	if err == srvorder.ErrOrderAlreadyTaken {
		log.WithField("err", err).Error("Failed to take order, since order already taken")
		ctx.StatusCode(iris.StatusConflict)
		ctx.JSON(iris.Map{
			"error": "Order already taken",
		})
		return
	}
	if err != nil {
		log.WithField("err", err).Error("Failed to take order")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"error": "Service unavailable",
		})
		return
	}

	log.Debug("Successfully took order")

	ctx.JSON(iris.Map{
		"status": "SUCCESS",
	})
}

func ListOrders(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "ListOrders"})

	page := ctx.Values().Get("_page").(int)
	limit := ctx.Values().Get("_limit").(int)
	offset := ctx.Values().Get("_offset").(int)

	log = log.WithFields(logrus.Fields{"page": page, "limit": limit, "offset": offset})

	orders, err := srvorder.ListOrders(offset, limit)
	if err != nil {
		log.WithField("err", err).Error("Failed to retrieve orders")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"error": "Failed to retrieve orders",
		})
		return
	}

	log.WithField("order_count", len(orders)).Debug("Successfully listed orders")
	ctx.JSON(orders)
}
