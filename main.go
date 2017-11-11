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
var count int

func main() {
	flag.StringVar(&host, "h", "", "Host to ping")
	flag.IntVar(&port, "p", 80, "Port to use for the TCP connection")
	flag.BoolVar(&debug, "d", false, "Output packet sent and received")
	flag.IntVar(&count, "c", 0, "Number of probes to send")
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

	if count > 0 {
		for i := 0; i < count; i++ {
			latency := probe.GetLatency(uint16(port))
			if latency > 0 {
				fmt.Printf("%s -> %s (%dms)\n", src, host, (int64(latency) / int64(100000)))
			} else {
				fmt.Printf("Timeout\n")
			}
		}
	} else {
		for {
			latency := probe.GetLatency(uint16(port))
			if latency > 0 {
				fmt.Printf("%s -> %s (%dms)\n", src, host, (int64(latency) / int64(100000)))
			} else {
				fmt.Printf("Timeout\n")
			}
			time.Sleep(time.Second)
		}
	}
}
