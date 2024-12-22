package main

import (
	"Clif/internal/filesystem"
	"Clif/internal/ui"
	"os"
)

func main() {
	initialDir, _ := os.Getwd()
	fileBrowser := filesystem.NewFileBrowser(initialDir)

	tui := ui.NewUI(fileBrowser)
	tui.Run()
}
