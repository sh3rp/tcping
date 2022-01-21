package tcp

import (
	"bytes"
	"log"
	"net"

	"inet.af/netaddr"
)

type TCPInjector interface {
	Inject(TCPPacket) error
}

type TCPReceiver interface {
	Receive(uint16) <-chan TCPPacket
}

func NewTCPInjectorWithConn(conn net.Conn) (TCPInjector, error) {
	return tcpInjector{conn, netaddr.IP{}}, nil
}

func NewTCPInjector(dstIP string) (TCPInjector, error) {
	conn, err := net.Dial("ip4:tcp", dstIP)

	if err != nil {
		return nil, err
	}

	return NewTCPInjectorWithConn(conn)
}

type tcpInjector struct {
	conn  net.Conn
	dstIP netaddr.IP
}

func (ti tcpInjector) Inject(pkt TCPPacket) error {
	if pkt.DstIP.IsZero() {
		pkt.DstIP = ti.dstIP
	}
	headerBytes := pkt.Header.MarshalTCP()
	pkt.Header.Checksum = Checksum(headerBytes, pkt.SrcIP.As4(), pkt.DstIP.As4())
	headerBytes = pkt.Header.MarshalTCP()
	buffer := bytes.Buffer{}
	buffer.Write(headerBytes)
	buffer.Write(pkt.Data)
	_, err := ti.conn.Write(buffer.Bytes())
	return err
}

func NewTCPReceiverWithConn(conn net.Conn) (TCPReceiver, error) {
	return tcpReceiver{conn, make(chan TCPPacket), uint16(conn.(*net.TCPConn).LocalAddr().(*net.TCPAddr).Port)}, nil
}

func NewTCPReceiver(srcIP string, listenPort uint16) (TCPReceiver, error) {
	netaddr, err := net.ResolveIPAddr("ip4", srcIP)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenIP("ip4:tcp", netaddr)
	if err != nil {
		return nil, err
	}
	return NewTCPReceiverWithConn(conn)
}

type tcpReceiver struct {
	conn          net.Conn
	packetChannel chan TCPPacket
	listenPort    uint16
}

func (tr tcpReceiver) Receive(port uint16) <-chan TCPPacket {
	go func() {
		var tcp *TCPHeader
		for {
			buf := make([]byte, 1460) // set to default MTU, should set this to currently configured MTU?
			_, raddr, err := tr.conn.(*net.IPConn).ReadFrom(buf)

			if err != nil {
				log.Printf("error reading packet: %+v", err)
				continue
			}

			tcp = ParseTCP(buf[:20])

			if tcp.DataOffset > 5 {
				for i := 21; i < len(buf[i:]); {
					tcp.WithOption(TCPOption{
						Kind:   buf[i],
						Length: buf[i+1],
						Data:   buf[i+2 : buf[i+1]],
					})
					i = i + 2 + int(buf[i+1])
				}
			}

			data := buf[20+(tcp.DataOffset*4):]

			tr.packetChannel <- NewTCPPacket("src", raddr.String(), tcp, data)
		}
	}()

	return tr.packetChannel
}
