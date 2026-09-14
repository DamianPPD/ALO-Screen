//go:build windows

package main

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	hotkeyID    = 1
	trayIconID  = 1
	menuCapture = 1001
	menuAbout   = 1002
	menuExit    = 1003
)

type app struct {
	hwnd       uintptr
	hotkeyOkay bool
	nid        notifyIconData
}

var currentApp *app
var wndProcCallback = syscall.NewCallback(wndProc)

func newApp() (*app, error) {
	// Używaj fizycznych pikseli na monitorach z różnym skalowaniem DPI.
	procSetProcessDpiAwarenessContext.Call(^uintptr(3)) // DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 (-4)

	hInstance, _, _ := procGetModuleHandleW.Call(0)
	if hInstance == 0 {
		return nil, errors.New("nie udało się pobrać uchwytu aplikacji")
	}

	className := utf16Ptr("ALO_SCREEN_HIDDEN_WINDOW")
	icon, _, _ := procLoadIconW.Call(0, idiApplication)
	cursor, _, _ := procLoadCursorW.Call(0, idcArrow)
	wc := wndClassEx{
		CbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		LpfnWndProc:   wndProcCallback,
		HInstance:     hInstance,
		HIcon:         icon,
		HCursor:       cursor,
		LpszClassName: className,
		HIconSm:       icon,
	}
	atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if atom == 0 {
		return nil, fmt.Errorf("nie udało się zarejestrować okna aplikacji: %v", err)
	}

	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(utf16Ptr("ALO Screen"))),
		0,
		0, 0, 0, 0,
		0, 0, hInstance, 0,
	)
	if hwnd == 0 {
		return nil, fmt.Errorf("nie udało się utworzyć ukrytego okna: %v", err)
	}

	a := &app{hwnd: hwnd}
	currentApp = a

	a.nid = notifyIconData{
		CbSize:           uint32(unsafe.Sizeof(notifyIconData{})),
		HWnd:             hwnd,
		UID:              trayIconID,
		UFlags:           nifMessage | nifIcon | nifTip,
		UCallbackMessage: wmTrayIcon,
		HIcon:            icon,
	}
	copyUTF16(a.nid.SzTip[:], "ALO Screen — kliknij, aby wkleić monitor 2")
	ret, _, err := procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&a.nid)))
	if ret == 0 {
		procDestroyWindow.Call(hwnd)
		return nil, fmt.Errorf("nie udało się dodać ikony obok zegara: %v", err)
	}

	versionData := a.nid
	versionData.UVersionTimeout = notifyIconVersion4
	procShellNotifyIconW.Call(nimSetVer, uintptr(unsafe.Pointer(&versionData)))

	ret, _, _ = procRegisterHotKey.Call(hwnd, hotkeyID, modNoRepeat, vkF8)
	a.hotkeyOkay = ret != 0
	return a, nil
}

func (a *app) close() {
	if a == nil {
		return
	}
	if a.hotkeyOkay {
		procUnregisterHotKey.Call(a.hwnd, hotkeyID)
	}
	procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&a.nid)))
	if a.hwnd != 0 {
		procDestroyWindow.Call(a.hwnd)
	}
	currentApp = nil
}

func (a *app) run() {
	var m msg
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(ret) <= 0 {
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func (a *app) captureAndPaste() {
	if err := captureTargetMonitorToClipboard(a.hwnd); err != nil {
		showError("ALO Screen — zrzut", err.Error())
		return
	}
	if err := activateChatGPTAndPaste(); err != nil {
		showError("ALO Screen — ChatGPT", err.Error()+"\n\nZrzut został już skopiowany do schowka — możesz wkleić go ręcznie Ctrl+V.")
	}
}

func (a *app) showMenu() {
	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer procDestroyMenu.Call(menu)

	procAppendMenuW.Call(menu, mfString, menuCapture, uintptr(unsafe.Pointer(utf16Ptr("Zrób zrzut monitora 2 i wklej"))))
	about := "ALO Screen 0.1 — F8: zrzut + wklejenie"
	if !a.hotkeyOkay {
		about = "ALO Screen 0.1 — F8 zajęty przez inny program"
	}
	procAppendMenuW.Call(menu, mfString, menuAbout, uintptr(unsafe.Pointer(utf16Ptr(about))))
	procAppendMenuW.Call(menu, mfString, menuExit, uintptr(unsafe.Pointer(utf16Ptr("Zamknij"))))

	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procSetForegroundWindow.Call(a.hwnd)
	procTrackPopupMenu.Call(menu, tpmRightButton|tpmBottomAlign|tpmLeftAlign, uintptr(int64(pt.X)), uintptr(int64(pt.Y)), 0, a.hwnd, 0)
}

func wndProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	a := currentApp
	switch message {
	case wmHotKey:
		if a != nil && wParam == hotkeyID {
			a.captureAndPaste()
			return 0
		}
	case wmTrayIcon:
		if a != nil {
			switch loword(lParam) {
			case wmLButtonUp:
				a.captureAndPaste()
				return 0
			case wmRButtonUp:
				a.showMenu()
				return 0
			}
		}
	case wmCommand:
		if a != nil {
			switch loword(wParam) {
			case menuCapture:
				a.captureAndPaste()
				return 0
			case menuAbout:
				text := "ALO Screen 0.1\n\nLewy klik ikony lub F8:\n1. zrzut całego monitora 2\n2. kopiowanie do schowka\n3. przełączenie do ChatGPT\n4. Ctrl+V\n\nEnter naciskasz sam."
				if !a.hotkeyOkay {
					text += "\n\nUwaga: F8 jest aktualnie zajęty przez inny program. Ikona w zasobniku nadal działa."
				}
				procMessageBoxW.Call(hwnd, uintptr(unsafe.Pointer(utf16Ptr(text))), uintptr(unsafe.Pointer(utf16Ptr("ALO Screen 0.1"))), mbOK)
				return 0
			case menuExit:
				procDestroyWindow.Call(hwnd)
				return 0
			}
		}
	case wmDestroy:
		if a != nil {
			if a.hotkeyOkay {
				procUnregisterHotKey.Call(hwnd, hotkeyID)
				a.hotkeyOkay = false
			}
			procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&a.nid)))
			a.hwnd = 0
		}
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(hwnd, uintptr(message), wParam, lParam)
	return ret
}
