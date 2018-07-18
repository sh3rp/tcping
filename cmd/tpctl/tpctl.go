package main

import (
	"context"

	"github.com/sh3rp/tcping/rpc"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080")
	if err != nil {
		rpc.LOGGER.Error("Error: %v", err)
	}
	defer conn.Close()

	client := rpc.NewTcpingServiceClient(conn)
	probes, err := client.GetProbes(context.Background(), &rpc.Empty{})
	rpc.LOGGER.Info("Probes: %v", probes)
}
