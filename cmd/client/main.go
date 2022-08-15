package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-training/proto-connect-demo/gen/go/gitea/v1"
	"github.com/go-training/proto-connect-demo/gen/go/gitea/v1/giteaconnect"
	"github.com/go-training/proto-connect-demo/gen/go/ping/v1"
	"github.com/go-training/proto-connect-demo/gen/go/ping/v1/pingconnect"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
)

func main() {
	c := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
	}

	connectGiteaClient := giteaconnect.NewGiteaServiceClient(
		c,
		"http://localhost:8080/",
	)

	grpcGiteaClient := giteaconnect.NewGiteaServiceClient(
		c,
		"http://localhost:8080/",
		connect.WithGRPC(),
	)

	giteaClients := []giteaconnect.GiteaServiceClient{connectGiteaClient, grpcGiteaClient}

	for _, client := range giteaClients {
		req := connect.NewRequest(&gitea.GiteaRequest{
			Name: "foobar",
		})
		req.Header().Set("Gitea-Header", "hello from connect")
		res, err := client.Gitea(context.Background(), req)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Message:", res.Msg.Giteaing)
		log.Println("Gitea-Version:", res.Header().Get("Gitea-Version"))
	}

	connectPingClient := pingconnect.NewPingServiceClient(
		c,
		"http://localhost:8080/",
	)

	grpcPingClient := pingconnect.NewPingServiceClient(
		c,
		"http://localhost:8080/",
		connect.WithGRPC(),
	)

	pingClients := []pingconnect.PingServiceClient{connectPingClient, grpcPingClient}

	for _, client := range pingClients {
		req := connect.NewRequest(&ping.PingRequest{
			Data: "Ping",
		})
		req.Header().Set("Gitea-Header", "hello from connect")
		res, err := client.Ping(context.Background(), req)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Message:", res.Msg.Data)
		log.Println("Gitea-Version:", res.Header().Get("Gitea-Version"))
	}
}
