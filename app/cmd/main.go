package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"webApp/configs"
	"webApp/controller"
	_ "webApp/docs"
	"webApp/usecase"
)

// @title Decision Maker APi
// @description RESTful API for program implementation MCDM methods TOPSIS and SMART

// @host localhost:3030
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	app := fiber.New(fiber.Config{})

	config := configs.SetConfig()

	di := usecase.NewDiContainer(config)
	defer func() {
		logrus.Fatal(di.ShutDown())
	}()

	hand := controller.NewHandler(di, &config.AppConfig)

	// Configure CORS for frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.FrontendURL,
		AllowHeaders:     "Authorization, Origin, Content-Type, Accept",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowCredentials: true,
	}))

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
