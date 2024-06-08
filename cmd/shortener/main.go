package main

import (
	"github.com/mikesvis/short/internal/app/config"
	"github.com/mikesvis/short/internal/app/server"
)

func main() {
	config.InitConfig()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
