package tcping

type ProbePacket struct {
	SrcIP string
	DstIP string
	Sent  *TCPHeader
	Recv  *TCPHeader
	Mark  float64
}

type Recvd struct {
	Recv *TCPHeader
	Mark float64
}
