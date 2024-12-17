package main

import (
	"net/http"

	"github.com/yyyoichi/submarine-game/internal/app"
)

func main() {
	mux := http.NewServeMux()
	app.Route(mux)
	http.ListenAndServe(":3000", mux)
}
