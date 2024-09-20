package main

import (
	_ "net/http/pprof"

	"github.com/mikesvis/short/internal/app"
)

/* Нужен 1 на 1 я не понял как снимать профиль памяти, но тесты прошли :) */
func main() {
	app := app.New()

	app.Run()
}
