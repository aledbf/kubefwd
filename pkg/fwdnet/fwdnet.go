package fwdnet

import (
	"fmt"
	"net"

	util "github.com/apcera/util/iprange"
	"github.com/sparrc/go-ping"
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

		pinger, err := ping.NewPinger(ip.String())
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		pinger.SetPrivileged(false)
		pinger.Run()

		stats := pinger.Statistics()
		if stats.PacketsRecv == 3 {
			// alive
			continue
		}

		return ip, nil
	}

	return nil, fmt.Errorf("No IP addresses available")
}
