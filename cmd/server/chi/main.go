package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/go-training/gin-connect-demo/gen/proto/v1"
	"github.com/go-training/gin-connect-demo/gen/proto/v1/v1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type GiteaServer struct{}

func (s *GiteaServer) Gitea(
	ctx context.Context,
	req *connect.Request[v1.GiteaRequest],
) (*connect.Response[v1.GiteaResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&v1.GiteaResponse{
		Giteaing: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}

func giteaHandler(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got connection:", r.Proto)
		h.ServeHTTP(w, r)
	})
}

func main() {
	greeter := &GiteaServer{}
	path, handler := v1connect.NewGiteaServiceHandler(greeter)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post(path+"{name}", giteaHandler(handler))
	http.ListenAndServe(
		":8080",
		h2c.NewHandler(r, &http2.Server{}),
	)
}
