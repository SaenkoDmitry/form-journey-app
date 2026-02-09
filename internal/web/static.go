package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var content embed.FS

func SPAHandler() http.Handler {
	sub, err := fs.Sub(content, "dist")
	if err != nil {
		panic(err)
	}
	fsHandler := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 1. API пути обрабатываются отдельно
		if strings.HasPrefix(path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// 2. Если файл реально существует, отдать
		if _, err := sub.Open(strings.TrimPrefix(path, "/")); err == nil {
			fsHandler.ServeHTTP(w, r)
			return
		}

		// 3. SPA fallback для React Router
		r.URL.Path = "/index.html"
		fsHandler.ServeHTTP(w, r)
	})
}
