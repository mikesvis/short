package main

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/mikesvis/short/internal/app"
	"github.com/mikesvis/short/internal/config"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func valueOrNA(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func buildInfo() string {
	return fmt.Sprintf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s",
		valueOrNA(buildVersion), valueOrNA(buildDate), valueOrNA(buildCommit),
	)
}

func main() {
	fmt.Println()

	config := config.NewConfig()
	app := app.New(config)

	app.Run()
}
