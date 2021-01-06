package server

import (
	"log"
	"net"

	pb "github.com/helloworlde/grpc-gateway/proto/api"
	"github.com/helloworlde/grpc-gateway/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var helloService = service.HelloService{}

func StartGrpcServer() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("Listen gRPC port failed: ", err)
	}

	server := grpc.NewServer()
	pb.RegisterHelloServiceServer(server, &helloService)
	// Health Check
	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.HealthCheck", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, healthServer)

	go RegisterToConsul("127.0.0.1:8500", 9090)

	log.Println("Start gRPC Server on 0.0.0.0:9090")
	err = server.Serve(listener)
	if err != nil {
		log.Fatalln("Start gRPC Server failed: ", err)
	}

}
