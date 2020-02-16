package main

import (
	"flag"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	// ignore remapped stop/cont signals sent by i3bar
	signal.Ignore(syscall.SIGUSR1, syscall.SIGUSR2)

	netIface := flag.String("net-iface", "", "")
	flag.Parse()

	logwriter, err := syslog.New(syslog.LOG_INFO, "bar3")
	if err == nil {
		log.SetOutput(logwriter)
	}

	slots := []Slot{
		VPN(time.Second),
		APT(30 * time.Minute),
	}
	for _, i := range strings.Split(*netIface, ",") {
		slots = append(slots, Bandwidth(i))
	}

	slots = append(slots,
		Music(),
		Volume(),
		RAM(time.Second),
		CPU(2*time.Second),
		Storage(10*time.Minute),
		Battery(time.Second),
		Weather(30*time.Minute),
		Date(),
	)

	w := NewWriter(os.Stdout)
	Run(w, Part{Text: "   ╱   ", Sat: .1, Lum: .3}, slots...)
}
