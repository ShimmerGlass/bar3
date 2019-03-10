package main

import (
	"bufio"
	"flag"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"os/signal"
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

	if _, ok := os.LookupEnv("BAR3_FORK"); !ok {
		log.Println("Startup: forking")

		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), "BAR3_FORK=1")
		cmd.Stdout = os.Stdout

		stderr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		go func() {
			r := bufio.NewScanner(stderr)
			for r.Scan() {
				log.Println("stderr: ", r.Text())
			}
		}()

		err = cmd.Run()
		if err != nil {
			log.Println("fork exited with error: ", err)
		}
		return
	}

	log.Println("Startup: fork")

	w := NewWriter(os.Stdout)
	Run(w, Style("  ╱  ", ColorInactive),
		VPN(time.Second),
		APT(30*time.Minute),
		Bandwidth(*netIface),
		Music(),
		Volume(),
		RAM(time.Second),
		CPU(2*time.Second),
		Storage(10*time.Minute),
		Battery(time.Minute),
		Weather(30*time.Minute),
		Date(),
	)
}
