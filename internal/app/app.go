package app

import (
	"SimpleVault/internal/config"
	"SimpleVault/internal/web"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func RunApp() {
	cfg := config.LoadConfig()

	tempEngine, err := web.NewTemplateEngine()
	if err != nil {
		panic(err)
	}
	app := fiber.New(fiber.Config{
		Views: tempEngine,
	})

	web.RegisterRoutes(app, cfg)

	if err := app.Listen(cfg.GetAddr()); err != nil {
		slog.Error("Start application with error", "detail", err.Error())
	}

}
