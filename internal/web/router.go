package web

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	sm := NewSessionManager()

	app.Get("/register", RegisterPage())
	app.Post("/register", Register())

	app.Get("/login", LoginPage(sm))
	app.Post("/login", Login(sm))

	protected := app.Group("/records", AuthRequired(sm))
	protected.Get("/", RecordsPage)
	protected.Get("/new", NewRecordPage())
	protected.Post("/new", CreateRecord())
	protected.Get("/search", SearchRecords())
	protected.Get("/gen-password", GeneratePasswordHandler())
	protected.Get("/:id", ViewRecord())
	protected.Get("/:id/download", DownloadFile())
	protected.Get("/:id/delete", DeleteRecordPage())
	protected.Post("/:id/delete", DeleteRecord())

	app.Get("/backup", AuthRequired(sm), BackupVault())
	app.Get("/restore", AuthRequired(sm), RestoreVaultPage())
	app.Post("/restore", AuthRequired(sm), RestoreVault())

	app.Post("/logout", AuthRequired(sm), Logout(sm))
}
