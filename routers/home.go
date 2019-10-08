package routers

import (
	hd "order-service/handlers"

	"github.com/kataras/iris"
)

func home(app *iris.Application) {
	app.Get("/", hd.Home)
}
