//go:build windows

package main

import "os"

func main() {
	app, err := newApp()
	if err != nil {
		showError("ALO Screen", err.Error())
		os.Exit(1)
	}
	defer app.close()
	app.run()
}
