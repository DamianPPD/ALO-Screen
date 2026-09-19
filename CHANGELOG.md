# Changelog

## 0.2.0 - 2026-09-19

Window-targeting fix.

- Removed the requirement for the active browser/window title to contain `ChatGPT`.
- Added foreground-window tracking.
- Added filtering for Windows shell/tray windows so tray clicks do not become the paste target.
- Added WinAPI support for `GetForegroundWindow` and `GetClassNameW`.
- Updated tray/about text to version 0.2.
- Updated the build script to produce `ALO-Screen-0.2.exe`.
- Kept the old ChatGPT-title lookup only as a fallback.
- Message submission remains manual; ALO Screen never presses Enter.

## 0.1.0 - 2026-09-14

Initial public source release.

- Added Windows system tray application.
- Added global `F8` capture hotkey.
- Added full-screen capture of monitor 2.
- Added Windows clipboard image support.
- Added automatic activation of an open ChatGPT window and `Ctrl+V` paste.
- Kept message submission manual; ALO Screen never presses Enter.
- Added tests for monitor selection and ChatGPT window-title matching.
