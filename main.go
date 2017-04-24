package main

import (
	"flag"
	"fmt"

	"github.com/sh3rp/tcping/tcping"
)

func main() {
	flag.Parse()
	host := flag.Arg(0)
	for {
		latency := tcping.GetLatency(tcping.GetInterface(), host, 80)
		fmt.Println("Host: %s (%dms)", host, latency)
	}
}
