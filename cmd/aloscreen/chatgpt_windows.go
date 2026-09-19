//go:build windows

package main

import (
	"errors"
	"syscall"
	"time"
	"unsafe"

	"aloscreen/internal/core"
)

func findChatGPTWindow() uintptr {
	var found uintptr
	cb := syscall.NewCallback(func(hwnd, lparam uintptr) uintptr {
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1
		}
		length, _, _ := procGetWindowTextLengthW.Call(hwnd)
		if length == 0 {
			return 1
		}
		buf := make([]uint16, int(length)+1)
		procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		title := syscall.UTF16ToString(buf)
		if core.IsChatGPTWindow(title) {
			found = hwnd
			return 0
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	return found
}

func windowClassName(hwnd uintptr) string {
	if hwnd == 0 {
		return ""
	}
	buf := make([]uint16, 256)
	n, _, _ := procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:n])
}

func isSystemShellWindow(hwnd uintptr) bool {
	switch windowClassName(hwnd) {
	case "Shell_TrayWnd", "Shell_SecondaryTrayWnd", "Progman", "WorkerW":
		return true
	default:
		return false
	}
}

func activateWindowAndPaste(hwnd uintptr) error {
	if hwnd == 0 {
		return errors.New("nie znaleziono okna, do którego można wkleić zrzut")
	}

	procShowWindow.Call(hwnd, swRestore)
	procBringWindowToTop.Call(hwnd)
	ret, _, _ := procSetForegroundWindow.Call(hwnd)
	if ret == 0 {
		return errors.New("Windows nie pozwolił przełączyć fokusu na poprzednio aktywne okno")
	}

	time.Sleep(180 * time.Millisecond)
	procKeybdEvent.Call(vkControl, 0, 0, 0)
	procKeybdEvent.Call(vkV, 0, 0, 0)
	procKeybdEvent.Call(vkV, 0, keyeventfKeyup, 0)
	procKeybdEvent.Call(vkControl, 0, keyeventfKeyup, 0)
	return nil
}
