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
var debug bool

func main() {
	flag.StringVar(&host, "h", "", "Host to ping")
	flag.IntVar(&port, "p", 80, "Port to use for the TCP connection")
	flag.BoolVar(&debug, "d", false, "Output packet sent and received")
	flag.Parse()

	if host == "" {
		fmt.Printf("Must supply a host with the -h option\n")
		os.Exit(1)
	}

	src := tcping.GetInterface()

	probe := tcping.NewProbe(src, host, debug)

	if debug {
		fmt.Printf("Src: %s\n", src)
	}

	for {
		latency := probe.GetLatency(uint16(port))
		if latency > 0 {
			fmt.Printf("%s -> %s (%dms)\n", src, host, (int64(latency) / int64(100000)))
		} else {
			fmt.Printf("Timeout")
		}
		time.Sleep(time.Second)
	}
}
