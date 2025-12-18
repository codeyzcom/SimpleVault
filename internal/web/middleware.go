package web

import "github.com/gofiber/fiber/v2"

const sessionCookie = "vault_session"

func AuthRequired(sm *SessionManager) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Cookies(sessionCookie)
		if id == "" {
			return ctx.Redirect("/login")
		}

		v, ok := sm.Get(id)
		if !ok {
			ctx.ClearCookie(sessionCookie)
			return ctx.Redirect("/login")
		}

		ctx.Locals("vault", v)
		return ctx.Next()
	}
}
