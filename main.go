package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sh3rp/tcping/tcping"
)

var host string
var port int

func main() {
	flag.StringVar(&host, "h", "", "Host to ping")
	flag.IntVar(&port, "p", 80, "Port to use for the TCP connection")
	flag.Parse()

	if host == "" {
		fmt.Printf("Must supply a host with the -h option\n")
		os.Exit(1)
	}

	src := tcping.GetInterface()

	fmt.Printf("Src: %s\n", src)

	for {
		latency := tcping.GetLatency(src, host, uint16(port))
		fmt.Printf("%s -> %s (%dms)\n", src, host, latency)
		time.Sleep(time.Second)
	}
}
