# ALO Screen

A small portable Windows utility for taking a full screenshot of monitor 2 and pasting it into the window you are currently using with one click or the global `F8` hotkey.

ALO Screen is designed for a simple ChatGPT workflow: keep the content you want to show on monitor 2, keep the ChatGPT message field active on monitor 1, press `F8`, review the pasted screenshot, and send the message yourself.

## Current version

**ALO Screen 0.2**

Version 0.2 fixes the window-detection problem from 0.1. The app no longer depends on the browser window title containing the word `ChatGPT`.

## Features

- Global `F8` hotkey
- Left-click the tray icon to capture
- Captures the entire second monitor
- Copies the screenshot directly to the Windows clipboard
- Remembers the foreground user window
- Returns to that window and sends `Ctrl+V`
- Does not require `ChatGPT` in the browser/window title
- **Does not press Enter** — you review and send the message yourself
- No installer required
- No autostart
- No ChatGPT API key required
- No network communication performed by ALO Screen itself

## Typical workflow

1. Open ChatGPT in Edge, Chrome, or another browser.
2. Put the content you want to show on monitor 2.
3. Click the ChatGPT message field so the ChatGPT window is active.
4. Press `F8`.
5. ALO Screen captures monitor 2, copies the image to the clipboard, and pastes it into the active ChatGPT message field.
6. Review the screenshot and press Enter yourself.

You can also left-click the tray icon. ALO Screen tracks the last normal foreground window so clicking the Windows tray does not normally lose the paste target.

## Tray menu

Right-click the tray icon to access:

- **Capture monitor 2 and paste**
- **ALO Screen 0.2 / F8 status**
- **Exit**

If Windows initially hides the icon under the `^` overflow menu, drag it next to the clock if you want it to remain visible.

## Requirements

- Windows 10 or Windows 11
- x64 system
- Two monitors

The current version prefers Windows display device `DISPLAY2`. If `DISPLAY2` is not present, it can fall back to the only non-primary display when exactly one secondary monitor exists.

## Important behavior

ALO Screen does not use browser automation and does not search for the ChatGPT input field. It pastes with normal Windows `Ctrl+V` into the selected foreground window.

For the most reliable ChatGPT workflow, click the message field before pressing `F8`.

If Windows refuses to restore the target window, the screenshot remains in the clipboard and can still be pasted manually with `Ctrl+V`.

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
go build -trimpath -ldflags="-H=windowsgui -s -w" -o ALO-Screen-0.2.exe ./cmd/aloscreen
```

Or use:

```text
build.bat
```

The build script creates `ALO-Screen-0.2.exe` and runs the test suite first.

## Project structure

```text
cmd/aloscreen/       Windows tray application, screenshot capture and clipboard logic
internal/core/       Platform-independent monitor/window selection logic
go.mod               Go module definition
build.bat            Windows build script
```

## Privacy

ALO Screen does not upload screenshots by itself and does not use the OpenAI API. It captures the selected monitor, places the image in the local Windows clipboard, switches back to the selected window, and simulates `Ctrl+V`.

The screenshot is sent only after you manually submit the message.

## Security note

Unofficial builds may trigger Microsoft SmartScreen because they are not digitally signed and may have little or no reputation. Review the source and build it yourself if you prefer.

## License

MIT License. See [LICENSE](LICENSE).

## Disclaimer

This is an independent utility and is not affiliated with or endorsed by OpenAI. ChatGPT is a trademark of OpenAI.

---

## Polski

**ALO Screen 0.2** to małe narzędzie portable dla Windows.

`F8` albo kliknięcie ikony przy zegarze:

1. wykonuje zrzut całego monitora 2,
2. kopiuje go do schowka,
3. wraca do aktywnego / ostatnio używanego okna,
4. wykonuje `Ctrl+V`.

Najważniejsza zmiana względem 0.1: program **nie wymaga już słowa „ChatGPT” w tytule okna**.

Najpewniejszy sposób użycia z ChatGPT: kliknij pole wiadomości, a następnie naciśnij `F8`.

Program **nie naciska Enter**. Użytkownik sam sprawdza zrzut i wysyła wiadomość.
