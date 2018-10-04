package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func ReadLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	return ret, nil
}

func readBandwidth(dev string) (int, int, error) {
	lines, err := ReadLines("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		if key != dev {
			continue
		}

		value := strings.Fields(strings.TrimSpace(fields[1]))

		r, err := strconv.Atoi(value[0])
		if err != nil {
			return 0, 0, err
		}

		t, err := strconv.Atoi(value[8])
		if err != nil {
			return 0, 0, err
		}

		return r, t, nil
	}

	return 0, 0, fmt.Errorf("dev \"%s\" not found", dev)
}

func Bandwidth(iface string) Slot {
	lastR, lastT := -1, -1
	return NewTimedSlot(time.Second, func() string {
		r, t, err := readBandwidth(iface)
		if err != nil {
			errMsg(err.Error())
		}

		if lastR == -1 && lastT == -1 {
			lastR = r
			lastT = t
			return ""
		}

		cr, ct := r-lastR, t-lastT
		lastR, lastT = r, t

		rColor, tColor := ColorInactive, ColorInactive
		if cr > 0 {
			rColor = ColorActive
		}
		if ct > 0 {
			tColor = ColorActive
		}

		return fmt.Sprintf(`<span font_size="small">%s</span> %s  %s <span font_size="small">%s</span>`,
			iconC("\uf063", rColor),
			humanize.Bytes(uint64(cr)),
			humanize.Bytes(uint64(ct)),
			iconC("\uf062", tColor),
		)
	})
}
