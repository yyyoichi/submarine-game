package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	apiv1 "github.com/yyyoichi/submarine-game/internal/gen/api/v1"
	"github.com/yyyoichi/submarine-game/internal/gen/api/v1/apiv1connect"
	"github.com/yyyoichi/submarine-game/internal/handler"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":8080"
	}

	rpc := http.NewServeMux()
	rpc.Handle(apiv1connect.NewHelloServiceHandler(&Handler{}))
	gmHandler := handler.NewHandler(context.Background())
	rpc.Handle(apiv1connect.NewGameServiceHandler(gmHandler))

	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)
	mux.Handle("/rpc/", http.StripPrefix("/rpc", rpc))

	if err := http.ListenAndServe(port, h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Panic(err)
	}
}

//go:embed all:web/dist
var assets embed.FS

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := tryRead(r.URL.Path, w)
	if err == nil {
		return
	}
	_ = tryRead("index.html", w)
}

func tryRead(requestedPath string, w http.ResponseWriter) error {
	reqPath := path.Join("web/dist", requestedPath)
	if reqPath == "web/dist" {
		reqPath = "web/dist/index.html"
	}

	if extension := strings.LastIndex(reqPath, "."); extension == -1 {
		reqPath = fmt.Sprintf("%s.html", reqPath)
	}

	// read file
	f, err := assets.Open(reqPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// dir check
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.ErrUnsupported
	}

	// content type check
	ext := filepath.Ext(requestedPath)
	var contentType string

	if m := mime.TypeByExtension(ext); m != "" {
		contentType = m
	} else {
		contentType = "text/html"
	}

	w.Header().Set("Content-Type", contentType)
	io.Copy(w, f)

	return nil
}

type Handler struct {
	apiv1connect.HelloServiceHandler
}

func (h *Handler) Say(ctx context.Context, req *connect.Request[apiv1.SayRequest]) (*connect.Response[apiv1.SayResponse], error) {
	slog.InfoContext(ctx, "get request", slog.String("name", req.Msg.Name))
	resp := &apiv1.SayResponse{
		Hello: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	}
	return connect.NewResponse(resp), nil
}

// func (h *handler) Count(ctx context.Context, req *connect.Request[apiv1.CountRequest], strem *connect.ServerStream[apiv1.CountResponse]) error {
// 	for i := range 100 {
// 		if err := strem.Send(&apiv1.CountResponse{Count: int64(i)}); err != nil {
// 			return err
// 		}
// 		select {
// 		case <-ctx.Done():
// 		case <-time.After(time.Duration(1) * time.Second):
// 		}
// 	}
// 	return nil
// }
