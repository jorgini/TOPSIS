package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"webApp/app/configs"
	"webApp/app/handlers"
	"webApp/app/usecase"
)

/* TODO
собрать всю бизнес логику и структуры в entity
*/

func main() {
	app := fiber.New(fiber.Config{})

	config := configs.SetConfig()

	di := usecase.NewDiContainer(config)
	defer func() {
		logrus.Fatal(di.ShutDown())
	}()

	hand := handlers.NewHandler(di, &config.AppConfig)

	app.Route("/", hand.SetAllRoutes)

	logrus.WithFields(logrus.Fields{
		"port": config.AppConfig.Port,
	}).Infoln("Starting a web-server on port")

	go func() {
		if err := app.Listen(config.AppConfig.Port); err != nil {
			logrus.Fatal(err)
		}
	}()

	logrus.Info("webApp successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("webApp shut down")

	if err := app.Shutdown(); err != nil {
		logrus.Errorf("error occurred on server shut down: %s", err.Error())
	}
}
