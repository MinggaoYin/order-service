package middlewares

import (
	"strconv"

	"order-service/models"
	srvorder "order-service/services"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func FetchOrder(service srvorder.OrderService) context.Handler {
	return func(ctx iris.Context) {
		idStr := ctx.Params().Get("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id < 0 {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"error": "order_id must be an integer",
			})
			return
		}

		order, err := service.GetById(id)
		if err == models.ErrNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{
				"error": "Order not found",
			})
			return
		}
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"error": "Failed to find order",
			})
			return
		}

		ctx.Values().Set("_order", order)
		ctx.Next()
	}
}
