//go:build windows

package nexodus

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/nexodus-io/nexodus/internal/util"
	"go.uber.org/zap"
)

// RouteExistsOS checks to see if a route exists for the specified prefix
func RouteExistsOS(prefix string) (bool, error) {
	if err := ValidateCIDR(prefix); err != nil {
		return false, err
	}

	var output bytes.Buffer
	var cmd *exec.Cmd

	if util.IsIPv4Prefix(prefix) {
		cmd = exec.Command("netsh", "interface", "ipv4", "show", "route")
		cmd.Stdout = &output
	}

	if util.IsIPv6Prefix(prefix) {
		cmd = exec.Command("netsh", "interface", "ipv6", "show", "route")
		cmd.Stdout = &output
	}

	if err := cmd.Run(); err != nil {
		return false, err
	}

	// Validate the IP we're expecting is in the output
	return strings.Contains(output.String(), prefix), nil
}

// AddRoute adds a windows route to the specified interface
func AddRoute(prefix, dev string) error {
	// TODO: replace with powershell
	_, err := RunCommand("netsh", "interface", "ipv4", "add", "route", prefix, dev)
	if err != nil {
		return fmt.Errorf("no windows route added: %w", err)
	}

	return nil
}

// discoverLinuxAddress only used for windows build purposes
func discoverLinuxAddress(logger *zap.SugaredLogger, family int) (net.IP, error) {
	return nil, nil
}

func linkExists(wgIface string) bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %w\n", err))
		return false
	}
	for _, iface := range ifaces {
		if iface.Name == wgIface {
			return true
		}
	}
	return false
}

// delLink only used for windows build purposes
func delLink(wgIface string) (err error) {
	return nil
}

// DeleteRoute deletes a windows route
func DeleteRoute(prefix, dev string) error {
	// netsh interface ip delete route [prefix] [interface|*] [gateway]
	_, err := RunCommand("netsh", "interface", "ipv4", "del", "route", prefix, dev)
	if err != nil {
		return fmt.Errorf("no route deleted: %w", err)
	}

	return nil
}

func defaultTunnelDevOS() string {
	return wgIface
}

// binaryChecks validate the required binaries are available
func binaryChecks() error {
	if !IsCommandAvailable(wgWinBinary) {
		return fmt.Errorf("%s command not found, is wireguard installed?", wgWinBinary)
	}
	return nil
}

// prepOS perform OS specific OS changes
func prepOS(logger *zap.SugaredLogger) error {
	// ensure the windows wireguard directory exists
	if err := CreateDirectory(WgWindowsConfPath); err != nil {
		return fmt.Errorf("unable to create the wireguard config directory [%s]: %w", WgWindowsConfPath, err)
	}
	return nil
}

// isIPv6Supported returns true if the platform supports IPv6
func isIPv6Supported() bool {
	// Use netsh to check IPv6 status on interfaces
	data, err := RunCommand("netsh", "interface", "ipv6", "show", "interfaces")
	if err != nil {
		return false
	}

	if strings.Contains(strings.ToLower(data), "disabled") {
		return false
	}

	return true
}
