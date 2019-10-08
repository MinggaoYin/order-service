package main

import (
	"order-service/routers"
	_ "order-service/startup"

	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	routers.Register(app)
	app.Run(iris.Addr(":8080"), iris.WithoutStartupLog)
}
