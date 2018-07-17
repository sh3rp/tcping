package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sh3rp/tcping/rpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	rpc.LOGGER.Info("Listening on port 8080")
	grpcServer := grpc.NewServer()
	rpc.LOGGER.Info("Starting server")
	rpc.RegisterTcpingServiceServer(grpcServer, &rpc.TcpingdServer{})
	grpcServer.Serve(lis)
}
