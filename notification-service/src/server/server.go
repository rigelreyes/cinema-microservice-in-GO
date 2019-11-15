package server

import (
	"cinemas-microservices/notification-service/src/api"
	"cinemas-microservices/notification-service/src/routes"
	"context"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var e *echo.Echo

// Start ...
func Start(r map[string]interface{}, se chan error) {

	e = echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	app := e.Group("/notification")

	routes.API(app, r["repo"].(api.Repository))
	routes.HealthyAPI(e)

	// Start server
	go func() {
		if err := e.Start(":" + strconv.Itoa(r["port"].(int))); err != nil {
			log.Info("shutting down the server")
			se <- err
		}
	}()
}

// Shutdown ...
func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Warn("Server Shutdown")
	os.Exit(1)
}

// GetServer ...
func GetServer() *echo.Echo {
	return e
}
