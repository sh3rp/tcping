package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sh3rp/tcping/tcping"
)

var VERSION = "1.3"

var host string
var port int
var iface string
var timeout int64
var debug bool
var count int
var showVersion bool

func main() {
	flag.IntVar(&port, "p", 80, "Port to use for the TCP connection")
	flag.BoolVar(&debug, "d", false, "Debug output packet sent and received")
	flag.IntVar(&count, "c", 0, "Number of probes to send")
	flag.StringVar(&iface, "i", "", "Interface to use as the source of the TCP packets")
	flag.Int64Var(&timeout, "t", 3000, "Time in milliseconds to wait for probe to return")
	flag.BoolVar(&showVersion, "v", false, "Version info")
	flag.Parse()

	if showVersion {
		fmt.Printf("tcping v%s\n", VERSION)
		return
	}

	host = flag.Arg(0)

	if host == "" {
		fmt.Printf("Must supply a host to ping.\n")
		os.Exit(1)
	}

	src := tcping.GetInterface(iface)

	probe := tcping.NewProbe(src, host, timeout, uint16(port), debug)

	if debug {
		fmt.Printf("Src IP: %s\n\n", src)
	}

	if count > 0 {
		for i := 0; i < count; i++ {
			sendProbe(probe, debug)
			time.Sleep(1 * time.Second)
		}
	} else {
		for {
			sendProbe(probe, debug)
			time.Sleep(1 * time.Second)
		}
	}
}

func sendProbe(probe tcping.Probe, debug bool) {
	result, err := probe.GetLatency()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if result.Latency() > 0 {
		if debug {
			fmt.Printf(tcping.FormatResult(result, true))
		} else {
			fmt.Printf("Sent from %-15s to %-15s: %d ms\n",
				probe.SrcIP,
				probe.DstIP,
				result.Latency()/int64(time.Millisecond))
		}
	} else {
		fmt.Printf("Timeout\n")
	}
}
