package web

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed dist/*
var content embed.FS

func SPAHandler() http.Handler {
	sub, _ := fs.Sub(content, "dist")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		if strings.HasPrefix(path, "api/") {
			http.NotFound(w, r)
			return
		}

		// Попробуем открыть файл
		if f, err := sub.Open(path); err == nil {
			defer f.Close()
			info, _ := f.Stat()
			if !info.IsDir() {
				// Читаем весь файл в память и создаём io.ReadSeeker
				data, _ := io.ReadAll(f)
				http.ServeContent(w, r, path, info.ModTime(), bytes.NewReader(data))
				return
			}
		}

		// SPA fallback
		fallback, _ := sub.Open("index.html")
		defer fallback.Close()
		info, _ := fallback.Stat()
		data, _ := io.ReadAll(fallback)
		http.ServeContent(w, r, "index.html", info.ModTime(), bytes.NewReader(data))
	})
}

// MountSPA подключает SPA на chi.Router, не ломая существующие api маршруты
func MountSPA(r chi.Router, pathPrefix string) {
	r.Handle(pathPrefix+"*", SPAHandler())
}
