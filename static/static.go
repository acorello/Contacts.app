package static

import (
	"embed"
	"net/http"
)

//go:embed *.js *.css
var static embed.FS

func FileServer() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.FS(static)))
}
