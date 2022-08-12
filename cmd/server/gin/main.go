package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	v1 "github.com/go-training/gin-connect-demo/gen/proto/v1"        // generated by protoc-gen-go
	"github.com/go-training/gin-connect-demo/gen/proto/v1/v1connect" // generated by protoc-gen-connect-go

	"github.com/bufbuild/connect-go"
	"github.com/gin-gonic/gin"
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

func giteaHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Got connection:", c.Request.Proto)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	greeter := &GiteaServer{}
	path, handler := v1connect.NewGiteaServiceHandler(greeter)

	r := gin.Default()
	r.UseH2C = true
	r.POST(path+":name", giteaHandler(handler))

	r.Run(":8080")
}