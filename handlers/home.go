package handlers

import (
	"github.com/kataras/iris"
)

func Home(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"status": "OK",
	})
}
