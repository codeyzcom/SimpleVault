package web

import (
	"SimpleVault/internal/config"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, cfg *config.Config) {
	sm := NewSessionManager(cfg.SessionTTL)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login", fiber.StatusPermanentRedirect)
	})

	app.Get("/register", RegisterPage())
	app.Post("/register", Register(cfg.DataStore))

	app.Get("/login", LoginPage(sm))
	app.Post("/login", Login(sm, cfg.DataStore))

	protected := app.Group("/records", AuthRequired(sm))
	protected.Get("/", RecordsPage)
	protected.Get("/new", NewRecordPage())
	protected.Post("/new", CreateRecord())
	protected.Get("/search", SearchRecords())
	protected.Get("/gen-password", GeneratePasswordHandler())
	protected.Get("/:id", ViewRecordPage())
	protected.Get("/:id/edit", EditRecordPage())
	protected.Post("/:id/edit", EditRecord())
	protected.Get("/:id/download", DownloadFile())
	protected.Get("/:id/delete", DeleteRecordPage())
	protected.Post("/:id/delete", DeleteRecord())

	app.Get("/backup", AuthRequired(sm), BackupVault())
	app.Get("/restore", AuthRequired(sm), RestoreVaultPage())
	app.Post("/restore", AuthRequired(sm), RestoreVault())

	app.Post("/logout", AuthRequired(sm), Logout(sm))
}
