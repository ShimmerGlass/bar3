package main

import (
	"flag"
	"log"
	"log/syslog"
	"os"
	"time"
)

func main() {
	netIface := flag.String("net-iface", "", "")
	flag.Parse()

	logwriter, err := syslog.New(syslog.LOG_INFO, "bar3")
	if err == nil {
		log.SetOutput(logwriter)
	}

	w := NewWriter(os.Stdout)
	Run(w, Style("  ╱  ", ColorInactive),
		Bandwidth(*netIface),
		Music(),
		Volume(),
		RAM(time.Second),
		CPU(2*time.Second),
		Date(),
	)
}
