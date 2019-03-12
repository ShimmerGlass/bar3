package main

import (
	"log"
	"time"

	"github.com/BellerophonMobile/gonetworkmanager"
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
		conns := nm.GetActiveConnections()
		for _, c := range conns {
			if c.GetVPN() {
				return []Part{IconPart("\uf456")}
			}
		}

		return nil
	})
}
