package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	v1 "github.com/go-training/proto-connect-demo/gen/proto/v1"
	"github.com/go-training/proto-connect-demo/gen/proto/v1/v1connect"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	compress1KB := connect.WithCompressMinBytes(1024)

	greeter := &GiteaServer{}
	path, handler := v1connect.NewGiteaServiceHandler(
		greeter,
		compress1KB,
	)

	// grpcV1
	grpcPath, gHandler := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(v1connect.GiteaServiceName),
		compress1KB,
	)

	// grpcV1Alpha
	grpcAlphaPath, gAlphaHandler := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(v1connect.GiteaServiceName),
		compress1KB,
	)

	// grpcHealthCheck
	grpcHealthPath, gHealthHandler := grpchealth.NewHandler(
		grpchealth.NewStaticChecker(v1connect.GiteaServiceName),
		compress1KB,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post(path+"{name}", giteaHandler(handler))
	r.Post(grpcPath+"{name}", giteaHandler(gHandler))
	r.Post(grpcAlphaPath+"{name}", giteaHandler(gAlphaHandler))
	r.Post(grpcHealthPath+"{name}", giteaHandler(gHealthHandler))
	http.ListenAndServe(
		":8080",
		h2c.NewHandler(r, &http2.Server{}),
	)
}
