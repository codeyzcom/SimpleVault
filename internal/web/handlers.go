package web

import (
	"SimpleVault/internal/crypto"
	"SimpleVault/internal/storage"
	"SimpleVault/internal/vault"
	"github.com/gofiber/fiber/v2"
	"io"
)

func RegisterPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{
			"Title": "Register",
		}, "layouts/public")
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
		return c.Render("login", fiber.Map{
			"Title": "Login",
			"Error": "Invalid password",
		}, "layouts/public")
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
			}, "layouts/private")
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
	}, "layouts/private")
}

func GeneratePasswordHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		p, _ := vault.GeneratePassword(16)
		return c.SendString(p)
	}
}

func NewRecordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("record_new", fiber.Map{
			"Title": "Add record",
		}, "layouts/private")
	}
}

func CreateRecord() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("vault").(*vault.VaultService)

		title := c.FormValue("title")

		switch c.FormValue("type") {
		case "note":
			return handleErr(c, v.AddNote(title, c.FormValue("text")))
		case "file":
			fh, err := c.FormFile("file")
			if err != nil {
				return err
			}
			f, err := fh.Open()
			defer f.Close()
			if err != nil {
				return err
			}
			data, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			return handleErr(c, v.AddFile(title, fh.Filename, data))
		case "credential":
			return handleErr(c, v.AddCredential(
				title,
				vault.CredentialData{

					Website:  c.FormValue("website"),
					Username: c.FormValue("username"),
					Password: c.FormValue("password"),
					Email:    c.FormValue("email"),
					Phone:    c.FormValue("phone"),
					Note:     c.FormValue("note"),
				}))
		default:
			return fiber.NewError(400, "unknown record type")
		}
	}
}

func ViewRecordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("vault").(*vault.VaultService)
		r, err := v.GetRecord(c.Params("id"))
		if err != nil {
			return fiber.ErrNotFound
		}

		return c.Render("record_view", fiber.Map{
			"Title":  "View record",
			"Record": r,
		}, "layouts/private")
	}
}

func DownloadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("vault").(*vault.VaultService)
		r, err := v.GetRecord(c.Params("id"))
		if err != nil || r.Type != vault.RecordFile {
			return fiber.ErrNotFound
		}
		c.Set("Content-Disposition", "attachment; filename="+r.File.Filename)
		return c.Send(r.File.Data)
	}
}

func DeleteRecordPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("record_delete", fiber.Map{
			"ID": c.Params("id"),
		}, "layouts/private")
	}
}

func DeleteRecord() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("vault").(*vault.VaultService)
		if err := v.DeleteRecord(c.Params("id")); err != nil {
			return err
		}
		return c.Redirect("/records")
	}
}

func SearchRecords() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("vault").(*vault.VaultService)
		q := c.Query("q")
		return c.Render("records", fiber.Map{
			"Records": v.Search(q),
			"Query":   q,
		}, "layouts/private")
	}
}

func BackupVault() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("vault").(*vault.VaultService)
		data, err := v.Export()
		if err != nil {
			return err
		}
		c.Set("Content-Disposition", "attachment; filename=vault.dat")
		return c.Send(data)
	}
}

func RestoreVaultPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(
			"restore", fiber.Map{"Title": "Restore"},
			"layouts/private")
	}
}

func RestoreVault() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("vault").(*vault.VaultService)
		fh, _ := c.FormFile("vault")
		f, _ := fh.Open()
		defer f.Close()

		data, _ := io.ReadAll(f)
		if err := v.Import(data); err != nil {
			return c.Render("restore", fiber.Map{
				"Error": err.Error(),
			})
		}
		return c.Redirect("/records")
	}
}

func handleErr(c *fiber.Ctx, err error) error {
	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return c.Redirect("/records")
}
