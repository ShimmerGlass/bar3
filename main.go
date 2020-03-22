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
	disks := flag.String("disks", "", "")
	flag.Parse()

	logwriter, err := syslog.New(syslog.LOG_INFO, "bar3")
	if err == nil {
		log.SetOutput(logwriter)
	}

	pw := &ProcessWatcher{}
	go pw.watch()

	slots := []Slot{
		Music(),
		Volume(),
		VPN(time.Second),
		APT(30 * time.Minute),
	}

	for _, i := range strings.Split(*netIface, ",") {
		slots = append(slots, Bandwidth(i))
	}

	slots = append(slots,
		RAM(time.Second, pw),
		CPU(2*time.Second, pw),
		Storage(10*time.Minute),
		DiskUtilisation(strings.Split(*disks, ","), time.Second),
		Battery(time.Second),
		Weather(30*time.Minute),
		Date(),
	)

	w := NewWriter(os.Stdout)
	Run(w, Part{Text: "  ╱  ", Sat: .1, Lum: .3}, slots...)
}
