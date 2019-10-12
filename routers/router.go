package routers

import (
	srvorder "order-service/services"

	"github.com/kataras/iris"
)

func Register(app *iris.Application, orderService srvorder.OrderService) {
	home(app)
	order(app, orderService)

	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)
}

func notFoundHandler(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"error": "Not found",
	})
}
