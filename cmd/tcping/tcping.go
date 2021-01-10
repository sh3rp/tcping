package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sh3rp/tcping"
)

var VERSION = "1.3.1"

var host string
var port int
var iface string
var timeout int64
var debug bool
var count int
var showVersion bool
var runAsService bool
var webservicePort int

func main() {
	flag.IntVar(&port, "p", 80, "Port to use for the TCP connection")
	flag.BoolVar(&debug, "d", false, "Debug output packet sent and received")
	flag.IntVar(&count, "c", 0, "Number of probes to send")
	flag.StringVar(&iface, "i", "", "Interface to use as the source of the TCP packets")
	flag.Int64Var(&timeout, "t", 3000, "Time in milliseconds to wait for probe to return")
	flag.BoolVar(&showVersion, "v", false, "Version info")
	flag.BoolVar(&runAsService, "S", false, "Run as a webservice in a webserver")
	flag.IntVar(&webservicePort, "P", 8080, "Port to run webservice on")
	flag.Parse()

	if showVersion {
		fmt.Printf("tcping v%s\n", VERSION)
		return
	}

	if runAsService {
		log.Printf("Starting webservice on port %d", webservicePort)
		src := tcping.GetInterface(iface)
		probe := tcping.NewProbe(src, time.Duration(time.Duration(timeout)*time.Millisecond), debug)
		log.Printf("%s", tcping.NewWebService(probe).Start(webservicePort))
		return
	}

	host = flag.Arg(0)

	if host == "" {
		fmt.Printf("Must supply a host to ping.\n")
		os.Exit(1)
	}

	src := tcping.GetInterface(iface)

	probe := tcping.NewProbe(src, time.Duration(time.Duration(timeout)*time.Millisecond), debug)

	if debug {
		fmt.Printf("Src IP: %s\n\n", src)
	}

	if count > 0 {
		for i := 0; i < count; i++ {
			sendProbe(probe, host, uint16(port))
			time.Sleep(1 * time.Second)
		}
	} else {
		for {
			sendProbe(probe, host, uint16(port))
			time.Sleep(1 * time.Second)
		}
	}
}

func sendProbe(probe tcping.Probe, dstIp string, dstPort uint16) {
	result, err := probe.GetLatency(dstIp, dstPort)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if result <= 0 {
		fmt.Printf("Sent from %-15s to %-15s: %d ms\n",
			probe.SrcIp,
			dstIp,
			result/int64(time.Millisecond))
	} else {
		fmt.Printf("Sent from %-15s to %-15s: timeout (%d ms)\n",
			probe.SrcIp,
			dstIp,
			result/1000)
	}
}
