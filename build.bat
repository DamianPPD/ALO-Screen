@echo off
setlocal
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go test ./...
if errorlevel 1 exit /b 1
go build -trimpath -ldflags="-H=windowsgui -s -w" -o ALO-Screen-0.2.exe ./cmd/aloscreen
if errorlevel 1 exit /b 1
echo Gotowe: ALO-Screen-0.2.exe
