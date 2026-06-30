// Package nas provides operations for controlling the NAS device:
//   - Wake-on-LAN to power on
//   - SSH poweroff to shut down
package nas

import (
	"fmt"
	"log/slog"
	"os/exec"
	"wol_admin/config"
)

// WOL sends a Wake-on-LAN magic packet to the NAS MAC address.
func WOL() error {
	mac := config.Cfg.NasMAC
	if mac == "" {
		return fmt.Errorf("nas_mac is not configured")
	}

	cmd := exec.Command("wakeonlan", mac)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("wakeonlan command failed", "mac", mac, "error", err, "output", string(output))
		return fmt.Errorf("wakeonlan failed: %w", err)
	}

	slog.Info("WOL packet sent", "mac", mac, "output", string(output))
	return nil
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
