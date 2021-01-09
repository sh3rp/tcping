package tcping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPHeaderCreateEmpty(t *testing.T) {
	tcp := NewTCPHeader()
	assert.Equal(t, tcp.Dst, uint16(0))
	assert.Equal(t, tcp.Src, uint16(0))
}

func TestTCPHeaderSetSrcDstPort(t *testing.T) {
	tcp := NewTCPHeader().SrcPort(100).DstPort(200)
	assert.Equal(t, tcp.Src, uint16(100))
	assert.Equal(t, tcp.Dst, uint16(200))
}

func TestTCPHeaderSetWindow(t *testing.T) {
	tcp := NewTCPHeader().Win(1000)
	assert.Equal(t, tcp.Window, uint16(1000))
}

func TestTCPHeaderSetSynAck(t *testing.T) {
	tcp := NewTCPHeader().WithFlag(SYN)
	assert.True(t, tcp.HasFlag(SYN))
	tcp = NewTCPHeader().WithFlag(ACK)
	assert.True(t, tcp.HasFlag(ACK))
}
