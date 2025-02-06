package app

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed public/*
var publicFiles embed.FS

func Route(mux *http.ServeMux) {
	api := http.NewServeMux()
	api.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hello, World!! %s", time.Now().Format(time.RFC3339))))
		w.WriteHeader(http.StatusOK)
	})
	public, err := fs.Sub(publicFiles, "public")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("ENV") == "local" {
		s := http.FileServer(http.Dir("/workspaces/submarine-game/internal/app/public"))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			s.ServeHTTP(w, r)
		})
	} else {
		mux.Handle("/", http.FileServer(http.FS(public)))
	}
	mux.Handle("/api/", api)
}
