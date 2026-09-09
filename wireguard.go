package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/vishvananda/netlink"
)

func GetWireGuardPeerIPs(interfaceName string) (map[string]string, error) {
	device, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("get interface %q: %w", interfaceName, err)
	}

	if _, ok := device.(*netlink.Wireguard); !ok {
		return nil, fmt.Errorf("interface %q is not a WireGuard interface", interfaceName)
	}

	output, err := exec.Command("wg", "show", interfaceName, "endpoints").Output()
	if err != nil {
		return nil, fmt.Errorf("read WireGuard endpoints: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	peerIPs := make(map[string]string, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[1] == "(none)" {
			continue
		}

		host, _, err := net.SplitHostPort(fields[1])
		if err != nil {
			continue
		}

		if net.ParseIP(host) == nil {
			continue
		}

		peerIPs[fields[0]] = host
	}

	return peerIPs, nil
}
