package public_assets

import (
	"embed"
	"net/http"
)

//go:embed vendored/*.js *.css
var _fs embed.FS

func FileServer() http.Handler {
	return http.FileServer(http.FS(_fs))
}
