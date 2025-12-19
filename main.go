package main

import (
	"SimpleVault/internal/app"
	"log/slog"
)

func main() {
	slog.Info(">>> Application starting")
	app.RunApp()
	slog.Info(">> Application shutdown")
}
