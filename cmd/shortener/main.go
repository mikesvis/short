package main

import (
	"github.com/mikesvis/short/internal/app"
	"github.com/mikesvis/short/internal/logger"
)

func main() {
	app := app.New()
	if err := logger.Initialize(); err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		logger.Log.Fatalw(err.Error(), "event", "start server")
	}
}
