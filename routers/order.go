package routers

import (
	hd "order-service/handlers"
	mid "order-service/middlewares"
	srvorder "order-service/services"

	"github.com/kataras/iris"
)

func order(app *iris.Application, orderService srvorder.OrderService) {
	app.Post("/orders", hd.PlaceOrder(orderService))
	app.Patch("/orders/:id", mid.FetchOrder(orderService), hd.TakeOrder(orderService))
	app.Get("/orders", mid.Paginate, hd.ListOrders(orderService))
}
