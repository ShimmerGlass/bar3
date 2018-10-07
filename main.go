package main

import (
	"flag"
	"os"
	"time"
)

func main() {
	netIface := flag.String("net-iface", "", "")
	flag.Parse()

	w := NewWriter(os.Stdout)
	Run(w, colorize(ColorInactive, "  ╱  "),
		Bandwidth(*netIface),
		Music(),
		Volume(),
		RAM(time.Second),
		CPU(2*time.Second),
		Date(),
	)
}
