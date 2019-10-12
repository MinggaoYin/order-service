package main

import (
	"os"

	"order-service/repositories"
	"order-service/routers"
	"order-service/services/distance"
	"order-service/services/order"
	"order-service/startup"

	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
)

func main() {
	startup.Init()

	orderRepo := repositories.NewMysqlOrderRepo(startup.Db)
	distanceCalculator, err := distance.NewDistanceService()
	if err != nil {
		log.WithError(err).Error("Failed to get distance service")
		os.Exit(1)
	}

	orderService := order.NewOrderService(orderRepo, distanceCalculator)

	app := iris.New()
	routers.Register(app, orderService)
	app.Run(iris.Addr(":8080"), iris.WithoutStartupLog)
}
