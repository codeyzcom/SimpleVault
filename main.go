package main

import (
	"SimpleVault/internal/web"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func main() {
	app := fiber.New(fiber.Config{
		Views:       web.NewTemplateEngine(),
		ViewsLayout: "layouts/main",
	})

	web.RegisterRoutes(app)

	if err := app.Listen(":7879"); err != nil {
		slog.Error("Start application with error", "detail", err.Error())
	}

	slog.Info("Application shutdown...")
}
