//go:build windows

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"aloscreen/internal/core"
)

func enumerateMonitors() ([]core.Monitor, error) {
	monitors := make([]core.Monitor, 0, 4)
	cb := syscall.NewCallback(func(hMonitor, hdcMonitor, lprcMonitor, dwData uintptr) uintptr {
		var info monitorInfoEx
		info.CbSize = uint32(unsafe.Sizeof(info))
		ret, _, _ := procGetMonitorInfoW.Call(hMonitor, uintptr(unsafe.Pointer(&info)))
		if ret == 0 {
			return 1
		}
		monitors = append(monitors, core.Monitor{
			Device:  syscall.UTF16ToString(info.SzDevice[:]),
			Primary: info.DwFlags&monitorinfofPrimary != 0,
			Left:    info.RcMonitor.Left,
			Top:     info.RcMonitor.Top,
			Right:   info.RcMonitor.Right,
			Bottom:  info.RcMonitor.Bottom,
		})
		return 1
	})
	ret, _, err := procEnumDisplayMonitors.Call(0, 0, cb, 0)
	if ret == 0 {
		return nil, fmt.Errorf("nie udało się odczytać monitorów: %v", err)
	}
	return monitors, nil
}

func captureTargetMonitorToClipboard(owner uintptr) error {
	monitors, err := enumerateMonitors()
	if err != nil {
		return err
	}
	target, ok := core.SelectTargetMonitor(monitors)
	if !ok {
		return errors.New("nie znaleziono monitora 2. Program wymaga dwóch monitorów")
	}
	return captureMonitorToClipboard(owner, target)
}

func captureMonitorToClipboard(owner uintptr, m core.Monitor) error {
	width := m.Right - m.Left
	height := m.Bottom - m.Top
	if width <= 0 || height <= 0 {
		return errors.New("monitor ma nieprawidłowe wymiary")
	}

	screenDC, _, _ := procGetDC.Call(0)
	if screenDC == 0 {
		return errors.New("nie udało się pobrać obrazu ekranu")
	}
	defer procReleaseDC.Call(0, screenDC)

	memDC, _, _ := procCreateCompatibleDC.Call(screenDC)
	if memDC == 0 {
		return errors.New("nie udało się utworzyć bufora obrazu")
	}
	defer procDeleteDC.Call(memDC)

	bmp, _, _ := procCreateCompatibleBitmap.Call(screenDC, uintptr(width), uintptr(height))
	if bmp == 0 {
		return errors.New("nie udało się utworzyć bitmapy zrzutu")
	}
	defer procDeleteObject.Call(bmp)

	old, _, _ := procSelectObject.Call(memDC, bmp)
	if old == 0 {
		return errors.New("nie udało się przygotować bitmapy zrzutu")
	}
	defer procSelectObject.Call(memDC, old)

	ret, _, _ := procBitBlt.Call(
		memDC,
		0,
		0,
		uintptr(width),
		uintptr(height),
		screenDC,
		uintptr(int64(m.Left)),
		uintptr(int64(m.Top)),
		srccopy|captureblt,
	)
	if ret == 0 {
		return errors.New("nie udało się wykonać zrzutu monitora")
	}

	pixelBytes := int(width) * int(height) * 4
	pixels := make([]byte, pixelBytes)
	bmi := bitmapInfo{}
	bmi.Header.BiSize = uint32(unsafe.Sizeof(bitmapInfoHeader{}))
	bmi.Header.BiWidth = width
	bmi.Header.BiHeight = height
	bmi.Header.BiPlanes = 1
	bmi.Header.BiBitCount = 32
	bmi.Header.BiCompression = biRGB
	bmi.Header.BiSizeImage = uint32(pixelBytes)

	lines, _, _ := procGetDIBits.Call(
		memDC,
		bmp,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&pixels[0])),
		uintptr(unsafe.Pointer(&bmi)),
		dibRGBColors,
	)
	if lines == 0 {
		return errors.New("nie udało się odczytać pikseli zrzutu")
	}

	return setDIBClipboard(owner, width, height, pixels)
}

func setDIBClipboard(owner uintptr, width, height int32, pixels []byte) error {
	const headerSize = 40
	total := headerSize + len(pixels)
	hmem, _, _ := procGlobalAlloc.Call(gmemMoveable, uintptr(total))
	if hmem == 0 {
		return errors.New("brak pamięci na obraz w schowku")
	}
	ownedByClipboard := false
	defer func() {
		if !ownedByClipboard {
			procGlobalFree.Call(hmem)
		}
	}()

	ptr, _, _ := procGlobalLock.Call(hmem)
	if ptr == 0 {
		return errors.New("nie udało się zablokować pamięci schowka")
	}
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint32(header[0:4], headerSize)
	binary.LittleEndian.PutUint32(header[4:8], uint32(width))
	binary.LittleEndian.PutUint32(header[8:12], uint32(height))
	binary.LittleEndian.PutUint16(header[12:14], 1)
	binary.LittleEndian.PutUint16(header[14:16], 32)
	binary.LittleEndian.PutUint32(header[16:20], biRGB)
	binary.LittleEndian.PutUint32(header[20:24], uint32(len(pixels)))
	procRtlMoveMemory.Call(ptr, uintptr(unsafe.Pointer(&header[0])), uintptr(len(header)))
	procRtlMoveMemory.Call(ptr+headerSize, uintptr(unsafe.Pointer(&pixels[0])), uintptr(len(pixels)))
	procGlobalUnlock.Call(hmem)

	var opened bool
	for i := 0; i < 12; i++ {
		ret, _, _ := procOpenClipboard.Call(owner)
		if ret != 0 {
			opened = true
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
	if !opened {
		return errors.New("schowek Windows jest zajęty")
	}
	defer procCloseClipboard.Call()

	ret, _, _ := procEmptyClipboard.Call()
	if ret == 0 {
		return errors.New("nie udało się wyczyścić schowka")
	}
	ret, _, _ = procSetClipboardData.Call(cfDIB, hmem)
	if ret == 0 {
		return errors.New("nie udało się umieścić zrzutu w schowku")
	}
	ownedByClipboard = true
	return nil
}
