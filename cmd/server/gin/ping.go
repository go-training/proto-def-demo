package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-training/proto-connect-demo/gen/go/ping/v1"
	"github.com/go-training/proto-connect-demo/gen/go/ping/v1/pingconnect"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/gin-gonic/gin"
)

type PingService struct{}

func (s *PingService) Ping(
	ctx context.Context,
	req *connect.Request[ping.PingRequest],
) (*connect.Response[ping.PingResponse], error) {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&ping.PingResponse{
		Data: fmt.Sprintf("Hello, %s!", req.Msg.Data),
	})
	res.Header().Set("Gitea-Version", "v1")
	return res, nil
}

func pingServiceRoute(r *gin.Engine) {
	compress1KB := connect.WithCompressMinBytes(1024)

	pingService := &PingService{}
	connectPath, connecthandler := pingconnect.NewPingServiceHandler(
		pingService,
		compress1KB,
	)

	// grpcV1
	grpcPath, gHandler := grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(pingconnect.PingServiceName),
		compress1KB,
	)

	// grpcV1Alpha
	grpcAlphaPath, gAlphaHandler := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(pingconnect.PingServiceName),
		compress1KB,
	)

	// grpcHealthCheck
	grpcHealthPath, gHealthHandler := grpchealth.NewHandler(
		grpchealth.NewStaticChecker(pingconnect.PingServiceName),
		compress1KB,
	)

	r.POST(connectPath+":name", grpcHandler(connecthandler))
	r.POST(grpcPath+"Ping", grpcHandler(gHandler))
	r.POST(grpcAlphaPath+"Ping", grpcHandler(gAlphaHandler))
	r.POST(grpcHealthPath+"Ping", grpcHandler(gHealthHandler))
}
