package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&v1.GiteaResponse{
		Giteaing: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}

func grpcHandler(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("protocol version:", r.Proto)
		h.ServeHTTP(w, r)
	})
}

func main() {
	compress1KB := connect.WithCompressMinBytes(1024)

	giteaService := &GiteaServer{}
	connectPath, connecthandler := v1connect.NewGiteaServiceHandler(
		giteaService,
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
	r.Post(connectPath+"{name}", grpcHandler(connecthandler))
	r.Post(grpcPath+"{name}", grpcHandler(gHandler))
	r.Post(grpcAlphaPath+"{name}", grpcHandler(gAlphaHandler))
	r.Post(grpcHealthPath+"{name}", grpcHandler(gHealthHandler))

	srv := &http.Server{
		Addr: ":8080",
		Handler: h2c.NewHandler(
			r,
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP listen and serve: %v", err)
	}
}
