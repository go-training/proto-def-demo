package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-training/proto-connect-demo/gen/go/gitea/v1"
	"github.com/go-training/proto-connect-demo/gen/go/gitea/v1/giteaconnect"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/gin-gonic/gin"
)

type GiteaServer struct{}

func (s *GiteaServer) Gitea(
	ctx context.Context,
	req *connect.Request[gitea.GiteaRequest],
) (*connect.Response[gitea.GiteaResponse], error) {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&gitea.GiteaResponse{
		Giteaing: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}

func giteaServiceRoute(r *gin.Engine) {
	compress1KB := connect.WithCompressMinBytes(1024)

	giteaService := &GiteaServer{}
	connectPath, connecthandler := giteaconnect.NewGiteaServiceHandler(
		giteaService,
		compress1KB,
	)

	// grpcV1
	grpcPath, gHandler := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(giteaconnect.GiteaServiceName),
		compress1KB,
	)

	// grpcV1Alpha
	grpcAlphaPath, gAlphaHandler := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(giteaconnect.GiteaServiceName),
		compress1KB,
	)

	// grpcHealthCheck
	grpcHealthPath, gHealthHandler := grpchealth.NewHandler(
		grpchealth.NewStaticChecker(giteaconnect.GiteaServiceName),
		compress1KB,
	)

	r.POST(connectPath+":name", grpcHandler(connecthandler))
	r.POST(grpcPath+"Gitea", grpcHandler(gHandler))
	r.POST(grpcAlphaPath+"Gitea", grpcHandler(gAlphaHandler))
	r.POST(grpcHealthPath+"Gitea", grpcHandler(gHealthHandler))
}
