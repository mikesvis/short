package main

import (
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/server"
)

func main() {
	if err := logger.Initialize(); err != nil {
		panic(err)
	}

	if err := server.Run(); err != nil {
		logger.Log.Fatalw(err.Error(), "event", "start server")
	}
}
