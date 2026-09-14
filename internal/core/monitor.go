package core

import "strings"

type Monitor struct {
	Device  string
	Primary bool
	Left    int32
	Top     int32
	Right   int32
	Bottom  int32
}

func SelectTargetMonitor(monitors []Monitor) (Monitor, bool) {
	for _, m := range monitors {
		if strings.EqualFold(m.Device, `\\.\DISPLAY2`) {
			return m, true
		}
	}

	var secondary Monitor
	count := 0
	for _, m := range monitors {
		if !m.Primary {
			secondary = m
			count++
		}
	}
	if count == 1 {
		return secondary, true
	}
	return Monitor{}, false
}
