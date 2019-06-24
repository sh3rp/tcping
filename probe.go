package tcping

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Probe struct {
	SrcIp   string
	Timeout int64 // in milliseconds
	debug   bool
	watcher ProbeWatcher
}

func NewProbe(srcIp string, timeout int64, debug bool) Probe {
	return Probe{
		SrcIp:   srcIp,
		Timeout: timeout,
		debug:   debug,
		watcher: NewProbeWatcher(srcIp),
	}
}

func (p Probe) GetLatency(dstIp string, dstPort uint16) (int64, error) {
	addrs, err := net.LookupHost(dstIp)
	if err != nil {
		log.Fatalf("Error resolving %s. %s\n", dstIp, err)
	}
	dstIp = addrs[0]
	var mark int64
	wg := &sync.WaitGroup{}
	wg.Add(1)
	p.watcher.WatchFor(dstIp, dstPort, func(tcp *TCPHeader) {
		if tcp.HasFlag(RST) || (tcp.HasFlag(SYN) && tcp.HasFlag(ACK)) {
			mark = time.Now().UnixNano()
			wg.Done()
		}
	})
	defer p.watcher.StopWatchFor(dstIp, dstPort)

	sendProbe, err := p.SendPing(p.SrcIp, dstIp, dstPort)

	isAlive := wait(wg, time.Duration(p.Timeout)*time.Millisecond)

	if isAlive {
		return (mark - sendProbe.Mark), nil
	} else {
		return 0, errors.New("TCPing to host timeout")
	}
}

func (p Probe) SendPing(srcIP, dstIP string, dstPort uint16) (ProbePacket, error) {

	tmpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", srcIP))

	if err != nil {
		return ProbePacket{}, err
	}

	l, err := net.ListenTCP("tcp", tmpAddr)
	if err != nil {
		return ProbePacket{}, err
	}
	defer l.Close()

	packet := NewTCPHeader().
		SrcPort(uint16(l.Addr().(*net.TCPAddr).Port)).
		DstPort(dstPort).
		SeqNum(rand.Uint32()).
		WithFlag(SYN).
		Win(32676)

	/*packet := TCPHeader{
		Src:        uint16(l.Addr().(*net.TCPAddr).Port),
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
	}*/

	data := packet.MarshalTCP()

	packet.Checksum = Checksum(data, to4byte(srcIP), to4byte(dstIP))

	data = packet.MarshalTCP()

	conn, err := net.Dial("ip4:tcp", dstIP)

	if err != nil {
		return ProbePacket{}, err
	}

	defer conn.Close()

	sendTime := time.Now().UnixNano()

	numWrote, err := conn.Write(data)

	if err != nil {
		return ProbePacket{}, err
	}

	if numWrote != len(data) {
		return ProbePacket{}, errors.New(fmt.Sprintf("Error writing %d/%d bytes\n", numWrote, len(data)))
	}

	return ProbePacket{srcIP, *packet, sendTime}, nil
}

// Grab first interface found and the first IP on it
func GetInterface(intf string) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error, no interfaces: %s\n", err)
		return ""
	}
	if intf != "" {
		for _, i := range interfaces {
			if i.Name == intf {
				addrs, err := i.Addrs()

				if err != nil {
					log.Printf(" %s. %s\n", i.Name, err)
					break
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
		}
	} else {
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

func printTCP(tcp *TCPHeader) {
	var str string
	str = fmt.Sprintf("[ SRC: %5d ] [ DST: %5d ]\n", tcp.Src, tcp.Dst)
	str = str + fmt.Sprintf("[ SEQ: %20d ]\n", tcp.Seq)
	str = str + fmt.Sprintf("[ ACK: %20d ]\n", tcp.Ack)
	str = str + fmt.Sprintf("[ FLG: ")
	if tcp.HasFlag(URG) {
		str = str + "U"
	} else {
		str = str + "_"
	}
	if tcp.HasFlag(ACK) {
		str = str + "A"
	} else {
		str = str + "_"
	}
	if tcp.HasFlag(PSH) {
		str = str + "P"
	} else {
		str = str + "_"
	}
	if tcp.HasFlag(RST) {
		str = str + "R"
	} else {
		str = str + "_"
	}
	if tcp.HasFlag(SYN) {
		str = str + "S"
	} else {
		str = str + "_"
	}
	if tcp.HasFlag(FIN) {
		str = str + "F"
	} else {
		str = str + "_"
	}
	str = str + fmt.Sprintf("]")
	str = str + fmt.Sprintf(" [ WIN: %5d ]\n", tcp.Window)
	str = str + fmt.Sprintf("[ SUM: %5d ] [ URG: %5d ] \n", tcp.Checksum, tcp.Urgent)
	for _, o := range tcp.Options {
		str = str + fmt.Sprintf("[ Option: kind=%d len=%d data=%v ]\n", o.Kind, o.Length, o.Data)
	}
	fmt.Printf(str)
}

func wait(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}
