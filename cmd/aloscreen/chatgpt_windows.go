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

func activateChatGPTAndPaste() error {
	hwnd := findChatGPTWindow()
	if hwnd == 0 {
		return errors.New("nie znaleziono otwartego okna ChatGPT. Otwórz ten czat w osobnym oknie Edge/Chrome i spróbuj ponownie")
	}

	procShowWindow.Call(hwnd, swRestore)
	procBringWindowToTop.Call(hwnd)
	ret, _, _ := procSetForegroundWindow.Call(hwnd)
	if ret == 0 {
		return errors.New("Windows nie pozwolił przełączyć fokusu na okno ChatGPT")
	}

	time.Sleep(160 * time.Millisecond)
	procKeybdEvent.Call(vkControl, 0, 0, 0)
	procKeybdEvent.Call(vkV, 0, 0, 0)
	procKeybdEvent.Call(vkV, 0, keyeventfKeyup, 0)
	procKeybdEvent.Call(vkControl, 0, keyeventfKeyup, 0)
	return nil
}
