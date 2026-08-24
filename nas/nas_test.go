package nas

import (
	"bytes"
	"net"
	"testing"
)

func TestBuildMagicPacket(t *testing.T) {
	mac := net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}

	packet := buildMagicPacket(mac)

	if len(packet) != 102 {
		t.Fatalf("magic packet length = %d, want 102", len(packet))
	}
	if !bytes.Equal(packet[:6], bytes.Repeat([]byte{0xFF}, 6)) {
		t.Fatalf("magic packet prefix = %x, want 6x FF", packet[:6])
	}
	for i := 0; i < 16; i++ {
		got := packet[6+i*6 : 6+(i+1)*6]
		if !bytes.Equal(got, mac) {
			t.Fatalf("mac block %d = %x, want %x", i, got, mac)
		}
	}
}
