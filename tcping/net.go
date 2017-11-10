package tcping

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetLatency(srcIP, dstIP string, dstPort uint16) time.Duration {
	var wg sync.WaitGroup
	wg.Add(1)
	var receiveTime time.Time

	addrs, err := net.LookupHost(dstIP)
	if err != nil {
		log.Fatalf("Error resolving %s. %s\n", dstIP, err)
	}
	dstIP = addrs[0]

	go func() {
		receiveTime = WaitForResponse(srcIP, dstIP, dstPort)
		wg.Done()
	}()

	time.Sleep(1 * time.Millisecond)
	sendTime := SendPing(srcIP, dstIP, 0, dstPort)

	wg.Wait()
	return receiveTime.Sub(sendTime)
}

func SendPing(srcIP, dstIP string, srcPort, dstPort uint16) time.Time {
	if srcPort == 0 {
		srcPort = getNextLocalPort()
	}

	packet := TCPHeader{
		Src:        srcPort,
		Dst:        dstPort,
		Seq:        rand.Uint32(),
		Ack:        0,
		DataOffset: 5,
		Reserved:   0,
		ECN:        0,
		Ctrl:       2,
		Window:     0xaaaa,
		Checksum:   0,
		Urgent:     0,
		Options:    []TCPOption{},
	}

	data := packet.MarshalTCP()

	packet.Checksum = Checksum(data, to4byte(srcIP), to4byte(dstIP))

	data = packet.MarshalTCP()

	conn, err := net.Dial("ip4:tcp", dstIP)
	if err != nil {
		log.Println("Dial: %s\n", err)
		return time.Now()
	}
	defer conn.Close()

	sendTime := time.Now()

	numWrote, err := conn.Write(data)

	if err != nil {
		log.Printf("Error writing: %v\n", err)
		return time.Now()
	}

	if numWrote != len(data) {
		log.Printf("Error writing %d/%d bytes\n", numWrote, len(data))
		return time.Now()
	}

	return sendTime
}

func WaitForResponse(localAddress, remoteAddress string, port uint16) time.Time {
	netaddr, err := net.ResolveIPAddr("ip4", localAddress)
	if err != nil {
		log.Printf("Error (resolve): net.ResolveIPAddr: %s. %s\n", localAddress, netaddr)
		return time.Now()
	}

	conn, err := net.ListenIP("ip4:tcp", netaddr)
	if err != nil {
		log.Printf("Error (listen): %s\n", err)
		return time.Now()
	}
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var receiveTime time.Time
	for {
		buf := make([]byte, 1024)
		numRead, raddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Error (read): %s\n", err)
			return time.Now()
		}
		if raddr.String() != remoteAddress {
			continue
		}
		receiveTime = time.Now()
		tcp := ParseTCP(buf[:numRead])
		if tcp.Dst == port && tcp.HasFlag(RST) {
			break
		}
	}
	return receiveTime
}

// Grab first interface found and the first IP on it
func GetInterface() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error, no interfaces: %s\n", err)
		return ""
	}
	for _, iface := range interfaces {
		if strings.HasPrefix(iface.Name, "lo") {
			continue
		}
		addrs, err := iface.Addrs()

		if err != nil {
			log.Printf(" %s. %s\n", iface.Name, err)
			continue
		}
		var retAddr net.Addr
		for _, a := range addrs {
			if !strings.Contains(a.String(), ":") {
				retAddr = a
				break
			}
		}
		if retAddr != nil {
			return retAddr.String()[:strings.Index(retAddr.String(), "/")]
		}
	}

	return ""
}

func to4byte(addr string) [4]byte {
	parts := strings.Split(addr, ".")
	b0, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("to4byte: %s (latency works with IPv4 addresses only, but not IPv6!)\n", err)
	}
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])
	return [4]byte{byte(b0), byte(b1), byte(b2), byte(b3)}
}

func getNextLocalPort() uint16 {
	return 0
}
