package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/neorestore/internal/api/restore"
)

var (
	sos logger.Logger
)

func Init() {
	sos := logger.NewConsoleLogger(logger.LogLevelDebug)
	sos.Infof("init restore-ec2")

}

func main() {
	// create an echo resource
	e := echo.New()

	// middleware : log every http request
	e.Use(middleware.Logger())

	// routes
	e.POST("/restore-prefix", restore.RestoreHandler)

	// start server
	port := 8080
	address := fmt.Sprintf(":%d", port)
	e.Logger.Fatal(e.Start(address))
}
