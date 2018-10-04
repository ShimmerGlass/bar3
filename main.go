package main

import (
	"os"
	"time"
)

func main() {
	w := NewWriter(os.Stdout)
	Run(w, color(ColorInactive, "  ╱  "),
		Bandwidth(),
		Music(),
		Volume(),
		RAM(time.Second),
		CPU(2*time.Second),
		Date(),
	)
}
