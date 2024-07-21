package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	apiv1 "github.com/yyyoichi/submarine-game/internal/gen/api/v1"
	"github.com/yyyoichi/submarine-game/internal/gen/api/v1/apiv1connect"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":8080"
	}

	rpc := http.NewServeMux()
	rpc.Handle(apiv1connect.NewHelloServiceHandler(&handler{}))

	mux := http.NewServeMux()
	// mux.HandleFunc("/", notFoundHandler)
	mux.Handle("/rpc/", http.StripPrefix("/rpc", rpc))

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Panic(err)
	}
}

type handler struct {
	apiv1connect.HelloServiceHandler
}

func (h *handler) Say(ctx context.Context, req *connect.Request[apiv1.SayRequest]) (*connect.Response[apiv1.SayResponse], error) {
	slog.InfoContext(ctx, "get request", slog.String("name", req.Msg.Name))
	resp := &apiv1.SayResponse{
		Hello: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	}
	return connect.NewResponse(resp), nil
}
