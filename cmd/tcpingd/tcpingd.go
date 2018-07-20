package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/sh3rp/tcping/rpc"
	"google.golang.org/grpc"
)

var port int

func main() {
	flag.IntVar(&port, "p", 8080, "Port to run server on")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	rpc.LOGGER.Info("Listening on port 8080")
	grpcServer := grpc.NewServer()
	rpc.LOGGER.Info("Starting server")
	rpc.RegisterTcpingServiceServer(grpcServer, &rpc.TcpingdServer{})
	grpcServer.Serve(lis)
}
