package main

import (
	"log"
	"time"

	"github.com/Wifx/gonetworkmanager"
)

func VPN(interval time.Duration) Slot {
	var nm gonetworkmanager.NetworkManager

	return NewTimedSlot(interval, func() []Part {
		if nm == nil {
			var err error
			nm, err = gonetworkmanager.NewNetworkManager()
			if err != nil {
				log.Println(err)
				return nil
			}
		}
		conns, err := nm.GetPropertyActiveConnections()
		if err != nil {
			log.Printf("vpn: %s", err)
			return nil
		}
		for _, c := range conns {
			isVPN, err := c.GetPropertyVPN()
			if err != nil {
				log.Printf("vpn: %s", err)
				continue
			}
			if isVPN {
				return []Part{IconPart("\uf456")}
			}
		}

		return nil
	})
}
