# ALO Screen

A tiny portable Windows utility for taking a screenshot of a second monitor and pasting it into an open ChatGPT conversation with one click or the `F8` hotkey.

ALO Screen is designed for a simple workflow: keep tools such as GitHub Desktop, PowerShell, a terminal, logs, or another application on monitor 2, then send a full-screen snapshot to ChatGPT without opening Snipping Tool, selecting an area, copying, and pasting manually.

## Features

- Global `F8` hotkey
- Left-click the tray icon to capture
- Captures the entire second monitor
- Copies the screenshot directly to the Windows clipboard
- Finds an open window whose title contains `ChatGPT`
- Brings that window to the foreground and sends `Ctrl+V`
- **Does not press Enter** — you review and send the message yourself
- No installer required
- No autostart
- No ChatGPT API key required
- No network communication performed by ALO Screen itself

## Typical workflow

1. Open ChatGPT in Edge, Chrome, or another browser window.
2. Click once in the ChatGPT message field so it has focus.
3. Keep the applications you want to show on monitor 2.
4. Press `F8` or left-click the ALO Screen tray icon.
5. ALO Screen captures monitor 2 and pastes the image into ChatGPT.
6. Review the screenshot and press Enter yourself.

## Tray menu

Right-click the tray icon to access:

- **Capture monitor 2 and paste**
- **ALO Screen 0.1 / F8 status**
- **Exit**

If Windows initially hides the icon under the `^` overflow menu, drag it next to the clock if you want it to remain visible.

## Requirements

- Windows 10 or Windows 11
- x64 system
- Two monitors
- An open ChatGPT window

The current version prefers Windows display device `DISPLAY2`. If `DISPLAY2` is not present, it can fall back to the only non-primary display when exactly one secondary monitor exists.

## Important limitation

ALO Screen activates the ChatGPT window and sends `Ctrl+V`, but it does not currently locate the message box by UI automation. The ChatGPT message field should already have focus. This is intentional in version 0.1 because it keeps the utility small and predictable.

If ChatGPT cannot be found or Windows refuses to change the foreground window, the screenshot remains in the clipboard and can still be pasted manually with `Ctrl+V`.

## Build from source

The project is written in Go and uses the Windows API directly.

Requirements for building:

- Go 1.23 or newer

Run the tests:

```powershell
go test ./...
```

Build a Windows GUI executable:

```powershell
go build -trimpath -ldflags="-H=windowsgui -s -w" -o ALO-Screen.exe ./cmd/aloscreen
```

Or use:

```text
build.bat
```

The build script creates `ALO-Screen-0.1.exe` and runs the test suite first.

## Project structure

```text
cmd/aloscreen/       Windows tray application, screenshot capture and clipboard logic
internal/core/       Platform-independent monitor/window selection logic
go.mod               Go module definition
build.bat            Windows build script
```

## Privacy

ALO Screen does not upload screenshots by itself and does not use the OpenAI API. It only captures the selected monitor, places the image in the local Windows clipboard, activates a ChatGPT window, and simulates `Ctrl+V`.

The screenshot is sent only after you manually submit the message in ChatGPT.

## Security note

Unofficial builds may trigger Microsoft SmartScreen because they are not digitally signed and may have little or no reputation. Review the source and build it yourself if you prefer.

## License

MIT License. See [LICENSE](LICENSE).

## Disclaimer

This is an independent utility and is not affiliated with or endorsed by OpenAI. ChatGPT is a trademark of OpenAI.

---

## Polski

ALO Screen to małe narzędzie portable dla Windows. `F8` albo kliknięcie ikony przy zegarze wykonuje zrzut całego monitora 2, kopiuje go do schowka, przełącza do otwartego okna ChatGPT i wykonuje `Ctrl+V`.

Program **nie naciska Enter**. Użytkownik sam sprawdza zrzut i wysyła wiadomość.

Wersja 0.1 nie używa API ChatGPT, nie wysyła nic samodzielnie do sieci i nie uruchamia się automatycznie z Windowsem.
