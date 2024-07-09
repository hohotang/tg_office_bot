package main

import (
	"log"
	"os"
	"tgbot/app"

	"golang.org/x/sys/windows"
)

var version string = "0"

// Constants for console modes
const (
	ENABLE_EXTENDED_FLAGS  = 0x0080
	ENABLE_QUICK_EDIT_MODE = 0x0040
)

func main() {
	tgApp := app.NewApp(version)
	tgApp.Run(true)
}

func init() {
	disableQuickEditMode()
}

// disableQuickEditMode disables the Quick Edit mode in the console.
func disableQuickEditMode() {
	hStdin := windows.Handle(os.Stdin.Fd())

	// Get current console mode
	var mode uint32
	err := windows.GetConsoleMode(hStdin, &mode)
	if err != nil {
		log.Printf("GetConsoleMode failed: %v\n", err)
		return
	}

	// Disable Quick Edit mode
	mode &^= ENABLE_QUICK_EDIT_MODE
	mode |= ENABLE_EXTENDED_FLAGS

	// Set new console mode
	err = windows.SetConsoleMode(hStdin, mode)
	if err != nil {
		log.Printf("SetConsoleMode failed: %v\n", err)
	}
}
