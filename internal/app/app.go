package app

import (
	"SimpleVault/internal/config"
	"SimpleVault/internal/web"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func RunApp() {
	cfg := config.LoadConfig()
	app := fiber.New(fiber.Config{
		Views: web.NewTemplateEngine(),
	})

	web.RegisterRoutes(app, cfg)

	if err := app.Listen(cfg.GetAddr()); err != nil {
		slog.Error("Start application with error", "detail", err.Error())
	}

}
