// Package nas provides operations for controlling the NAS device:
//   - Wake-on-LAN to power on
//   - SSH poweroff to shut down
package nas

import (
	"fmt"
	"log/slog"
	"net"
	"os/exec"
	"wol-panel/config"
)

// WOL builds and sends a Wake-on-LAN magic packet to the NAS MAC address.
//
// A magic packet is 6x 0xFF followed by the 6-byte MAC repeated 16 times
// (102 bytes total), broadcast over UDP to a broadcast address on port 9.
func WOL() error {
	mac := config.Cfg.NasMAC
	if mac == "" {
		return fmt.Errorf("nas_mac is not configured")
	}

	hw, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("invalid nas_mac %q: %w", mac, err)
	}
	if len(hw) != 6 {
		return fmt.Errorf("nas_mac %q is not a 6-byte hardware address", mac)
	}

	addr, err := net.ResolveUDPAddr("udp4", config.Cfg.WOLBroadcast)
	if err != nil {
		return fmt.Errorf("invalid wol_broadcast %q: %w", config.Cfg.WOLBroadcast, err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("open udp socket: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(buildMagicPacket(hw)); err != nil {
		slog.Error("WOL packet send failed", "mac", mac, "broadcast", addr.String(), "error", err)
		return fmt.Errorf("WOL failed: %w", err)
	}

	slog.Info("WOL packet sent", "mac", mac, "broadcast", addr.String())
	return nil
}

// buildMagicPacket assembles the 102-byte Wake-on-LAN magic packet.
func buildMagicPacket(mac net.HardwareAddr) []byte {
	packet := make([]byte, 0, 6+16*len(mac))
	for i := 0; i < 6; i++ {
		packet = append(packet, 0xFF)
	}
	for i := 0; i < 16; i++ {
		packet = append(packet, mac...)
	}
	return packet
}

// Shutdown sends an SSH poweroff command to the NAS.
func Shutdown() error {
	ip := config.Cfg.NasIP
	user := config.Cfg.NasUser
	if ip == "" || user == "" {
		return fmt.Errorf("nas_ip or nas_user is not configured")
	}

	target := user + "@" + ip
	cmd := exec.Command("ssh", target, "sudo poweroff")
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("ssh poweroff command failed", "target", target, "error", err, "output", string(output))
		return fmt.Errorf("ssh poweroff failed: %w", err)
	}

	slog.Info("shutdown command sent", "target", target)
	return nil
}
