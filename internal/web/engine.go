package web

import (
	"github.com/gofiber/template/html/v2"
)

func NewTemplateEngine() *html.Engine {
	return html.New("./internal/web/templates", ".html")
}
