package routers

import (
	"github.com/kataras/iris"
)

func Register(app *iris.Application) {
	home(app)
	order(app)

	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)
}

func notFoundHandler(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"error": "Not found",
	})
}
