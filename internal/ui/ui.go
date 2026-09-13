package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func Handler() http.Handler {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic("ui: failed to create sub FS: " + err.Error())
	}
	return http.FileServer(http.FS(sub))
}
