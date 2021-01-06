package server

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/helloworlde/grpc-gateway/proto/api"
	"github.com/simplesurance/grpcconsulresolver/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(consul.NewBuilder())
}

func StartGwServer() {
	conn, err := grpc.DialContext(
		context.Background(),
		"consul://127.0.0.1:8500/server?health=healthy",
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatalln("Failed to dial server: ", err)
	}

	mux := runtime.NewServeMux()
	err = pb.RegisterHelloServiceHandler(context.Background(), mux, conn)

	if err != nil {
		log.Fatalln("Failed to register gateway: ", err)
	}

	server := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Println("Start gRPC Gateway Server on http://0.0.0.0:8090")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("Start Gateway Server failed: ", err)
	}

}
