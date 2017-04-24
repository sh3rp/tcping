package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sh3rp/tcping/tcping"
)

func main() {
	flag.Parse()
	dst := flag.Arg(0)
	src := tcping.GetInterface()
	for {
		latency := tcping.GetLatency(src, dst, 80)
		fmt.Printf("%s -> %s (%dms)\n", src, dst, latency)
		time.Sleep(time.Second)
	}
}
