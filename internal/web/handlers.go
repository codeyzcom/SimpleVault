package web

import (
	"SimpleVault/internal/crypto"
	"SimpleVault/internal/storage"
	"SimpleVault/internal/vault"
	"github.com/gofiber/fiber/v2"
)

func RegisterPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{
			"Title": "Register",
		})
	}
}

func Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		v := vault.NewVaultService(
			crypto.NewCryptoService(),
			storage.NewFileStorage("data/"+username),
		)

		if err := v.Init(password); err != nil {
			return c.Render("register", fiber.Map{
				"Error": err.Error(),
			})
		}

		return c.Redirect("/login")
	}
}

func LoginPage(sm *SessionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Cookies(sessionCookie)
		if id != "" {
			if _, ok := sm.Get(id); ok {
				return c.Redirect("/records")
			}
		}
		return c.Render("login", nil)
	}
}

func Login(sm *SessionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		v := vault.NewVaultService(
			crypto.NewCryptoService(),
			storage.NewFileStorage("data/"+username),
		)

		if err := v.Login(password); err != nil {
			return c.Render("login", fiber.Map{
				"Error": "Invalid credentials",
			})
		}
		sid := sm.Create(username, v)

		c.Cookie(&fiber.Cookie{
			Name:     sessionCookie,
			Value:    sid,
			HTTPOnly: true,
			SameSite: "Strict",
		})

		return c.Redirect("/records")
	}
}

func Logout(sm *SessionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sm.Delete(c.Cookies(sessionCookie))
		c.ClearCookie(sessionCookie)
		return c.Redirect("/login")
	}
}

func RecordsPage(c *fiber.Ctx) error {
	v := c.Locals("vault").(*vault.VaultService)
	return c.Render("records", fiber.Map{
		"Records": v.List(),
	})
}

func AddRecord(c *fiber.Ctx) error {
	v := c.Locals("vault").(*vault.VaultService)
	if err := v.Add(c.FormValue("title"), c.FormValue("content")); err != nil {
		return err
	}
	return c.Redirect("/records")
}
