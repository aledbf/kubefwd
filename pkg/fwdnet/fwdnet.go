package fwdnet

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	util "github.com/apcera/util/iprange"
)

var (
	allocated = make(map[string]bool, 1)
)

func Allocate(iprange string) (net.IP, error) {
	ipr, err := util.ParseIPRange(iprange)
	if err != nil {
		return nil, err
	}

	ipAllocator := util.NewAllocator(ipr)

	var ip net.IP
	for {
		if ipAllocator.Remaining() == 0 {
			break
		}

		ip = ipAllocator.Allocate()
		if _, ok := allocated[ip.String()]; ok {
			continue
		}

		allocated[ip.String()] = true

		out, _ := exec.Command("ping", ip.String(), "-c 2", "-w 10").Output()
		if strings.Contains(string(out), "Destination Host Unreachable") {
			return ip, nil
		}
	}

	return nil, fmt.Errorf("No IP addresses available")
}
