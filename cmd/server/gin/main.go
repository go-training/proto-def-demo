package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	v1 "github.com/go-training/proto-connect-demo/gen/proto/v1"        // generated by protoc-gen-go
	"github.com/go-training/proto-connect-demo/gen/proto/v1/v1connect" // generated by protoc-gen-connect-go
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/gin-gonic/gin"
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

func giteaHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("protocol version:", c.Request.Proto)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	compress1KB := connect.WithCompressMinBytes(1024)

	greeter := &GiteaServer{}
	connectPath, connecthandler := v1connect.NewGiteaServiceHandler(
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

	r := gin.Default()
	r.POST(connectPath+":name", giteaHandler(connecthandler))
	r.POST(grpcPath+":name", giteaHandler(gHandler))
	r.POST(grpcAlphaPath+":name", giteaHandler(gAlphaHandler))
	r.POST(grpcHealthPath+":name", giteaHandler(gHealthHandler))

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
