package core

import "testing"

func TestSelectTargetMonitorPrefersDisplay2(t *testing.T) {
	monitors := []Monitor{
		{Device: `\\.\DISPLAY1`, Primary: true},
		{Device: `\\.\DISPLAY3`},
		{Device: `\\.\DISPLAY2`},
	}
	got, ok := SelectTargetMonitor(monitors)
	if !ok {
		t.Fatal("expected a target monitor")
	}
	if got.Device != `\\.\DISPLAY2` {
		t.Fatalf("expected DISPLAY2, got %q", got.Device)
	}
}

func TestSelectTargetMonitorFallsBackToOnlySecondary(t *testing.T) {
	monitors := []Monitor{
		{Device: `\\.\DISPLAY1`, Primary: true},
		{Device: `\\.\DISPLAY7`},
	}
	got, ok := SelectTargetMonitor(monitors)
	if !ok {
		t.Fatal("expected a target monitor")
	}
	if got.Device != `\\.\DISPLAY7` {
		t.Fatalf("expected sole secondary display, got %q", got.Device)
	}
}

func TestSelectTargetMonitorFailsWithOneDisplay(t *testing.T) {
	monitors := []Monitor{{Device: `\\.\DISPLAY1`, Primary: true}}
	if _, ok := SelectTargetMonitor(monitors); ok {
		t.Fatal("expected no target monitor when only one display exists")
	}
}
