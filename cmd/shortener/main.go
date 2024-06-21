package main

import (
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/server"
)

func main() {
	config.InitConfig()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
