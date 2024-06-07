package main

import (
	"github.com/mikesvis/short/internal/app/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
