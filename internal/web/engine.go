package web

import (
	"embed"
	"github.com/gofiber/template/html/v2"
	"io/fs"
	"net/http"
)

//go:embed templates/**/*.html templates/*.html
var templatesFS embed.FS

func NewTemplateEngine() (*html.Engine, error) {
	sub, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		return nil, err
	}
	return html.NewFileSystem(http.FS(sub), ".html"), nil
}
