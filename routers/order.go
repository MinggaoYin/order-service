package routers

import (
	hd "order-service/handlers"
	mid "order-service/middlewares"

	"github.com/kataras/iris"
)

func order(app *iris.Application) {
	app.Post("/orders", hd.PlaceOrder)
	app.Patch("/orders/:id", mid.FetchOrder, hd.TakeOrder)
	app.Get("/orders", mid.Paginate, hd.ListOrders)
}
