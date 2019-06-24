package tcping

import (
	"net"
)

type ProbeWatcher interface {
	WatchFor(string, uint16, func(*TCPHeader))
	StopWatchFor(string, uint16)
}

type probeWatcher struct {
	localAddress string
	watches map[string]map[uint16]func(*TCPHeader)
}

func NewProbeWatcher(localAddress string) ProbeWatcher{
	return probeWatcher{
		localAddress: localAddress,
		watches: make(map[string]map[uint16]func(*TCPHeader)),
	}
}

func (pw probeWatcher) WatchFor(srcIp string, srcPort uint16, f func(*TCPHeader)) {
	if _,ok := pw.watches[srcIp]; !ok {
		pw.watches[srcIp] = make(map[uint16]func(*TCPHeader))
	}
	pw.watches[srcIp][srcPort] = f
}

func (pw probeWatcher) StopWatchFor(srcIp string, srcPort uint16) {
	if _, ipExists := pw.watches[srcIp]; ipExists {
		if _, portExists := pw.watches[srcIp][srcPort]; portExists {
			delete(pw.watches[srcIp],srcPort)
		}
	}
}

func (pw probeWatcher) watch() {
	netaddr, err := net.ResolveIPAddr("ip4", pw.localAddress)
	if err != nil {
		return
	}

	conn, err := net.ListenIP("ip4:tcp", netaddr)
	if err != nil {
		return
	}

	var tcp *TCPHeader
	for {
		buf := make([]byte, 1024)
		numRead, raddr, err := conn.ReadFrom(buf)
		if err != nil {
			return
		}
		if _, watched := pw.watches[raddr.String()]; !watched {
			continue
		}
		tcp = ParseTCP(buf[:numRead])

		if f, hasPort := pw.watches[raddr.String()][tcp.Src]; hasPort {
			f(tcp)
		}
	}
}