package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"embed"
)

//go:embed public/*
var publicFiles embed.FS

func main() {
	mux := http.NewServeMux()
	api := http.NewServeMux()
	api.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hello, World!! %s", time.Now().Format(time.RFC3339))))
		w.WriteHeader(http.StatusOK)
	})
	public, err := fs.Sub(publicFiles, "public")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(public)))
	mux.Handle("/api/", api)
	http.ListenAndServe(":3000", mux)
}
