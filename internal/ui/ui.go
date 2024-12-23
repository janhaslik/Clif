package ui

import (
	"Clif/internal/filesystem"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	fileBrowser *filesystem.FileBrowser
	app         *tview.Application
	searchInput *tview.InputField
	statusBox   *tview.TextView
}

func NewUI(fileBrowser *filesystem.FileBrowser) *UI {
	return &UI{
		fileBrowser: fileBrowser,
		app:         tview.NewApplication(),
		searchInput: tview.NewInputField(),
		statusBox:   tview.NewTextView().SetDynamicColors(true),
	}
}

func (ui *UI) Run() {
	ui.displayDirectory()
	err := ui.app.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func (ui *UI) displayDirectory() {
	files := ui.fileBrowser.DirEntriesFilter()

	if len(files) == 0 {
		files = []os.DirEntry{}
	}

	list := tview.NewList().ShowSecondaryText(false)

	for _, file := range files {
		if file.IsDir() {
			list.AddItem(file.Name(), "", '>', func() {
				err := ui.fileBrowser.NavigateInto(file.Name())
				if err != nil {
					return
				}
				ui.displayDirectory()
			})
		} else {
			list.AddItem(file.Name(), "", '>', func() {
				ui.displayFile(file.Name())
			})
		}
	}

	list.AddItem("..", "", 'u', func() {
		err := ui.fileBrowser.NavigateUp()
		if err != nil {
			return
		}
		ui.displayDirectory()
	})

	list.AddItem("Quit", "", 'q', func() {
		ui.app.Stop()
		os.Exit(0)
	})

	ui.searchInput.SetLabel("Search: ").
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				ui.app.SetFocus(list)
			} else if key == tcell.KeyEnter {
				ui.fileBrowser.Search(ui.searchInput.GetText())
				ui.displayDirectory()
			}
		})

	ui.searchInput.
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorGreen).
		SetLabelStyle(tcell.StyleDefault.Bold(true))

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'f':
			ui.app.SetFocus(ui.searchInput)
			return nil
		}
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.searchInput, 1, 0, true).
		AddItem(list, 0, 1, false)

	ui.app.SetRoot(flex, true).SetFocus(list)

	ui.updateStatusBox(true)
}

func (ui *UI) displayFile(name string) {
	content, _ := ui.fileBrowser.GetFileContent(name)

	textArea := tview.NewTextArea().
		SetText(content, false)

	isSaved := true
	ui.updateStatusBox(true)

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			err := ui.fileBrowser.SaveFileContent(name, textArea.GetText())
			if err != nil {
				return nil
			}
			isSaved = true
			ui.updateStatusBox(isSaved)
		case tcell.KeyEscape:
			ui.handleEscapeKey(name, textArea, &isSaved)
			return nil
		}
		isSaved = content == textArea.GetText()
		ui.updateStatusBox(isSaved)
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textArea, 0, 1, true).
		AddItem(ui.statusBox, 1, 0, false)

	ui.app.SetRoot(flex, true).SetFocus(textArea)

	err := ui.app.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func (ui *UI) handleEscapeKey(name string, textArea *tview.TextArea, isSaved *bool) {
	if *isSaved {
		ui.displayDirectory()
	} else {
		ui.statusBox.SetText("[yellow::b]You have unsaved changes. Press [white::b]y[yellow::b] to save, [white::b]n[yellow::b] to cancel, [white::b]q[yellow::b] to quit without saving[white]")

		textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'y':
				err := ui.fileBrowser.SaveFileContent(name, textArea.GetText())
				if err != nil {
					log.Printf("Error saving file: %v", err)
					return nil
				}
				ui.displayDirectory()
			case 'n':
				ui.app.SetFocus(textArea)
				ui.updateStatusBox(false)
			case 'q':
				ui.displayDirectory()
			}
			return nil
		})
	}
}

func (ui *UI) updateStatusBox(isSaved bool) {
	if isSaved {
		ui.statusBox.SetText("[green::b]Saved[white]").SetBackgroundColor(tcell.ColorBlack)
	} else {
		ui.statusBox.SetText("[red::b]Unsaved[white]").SetBackgroundColor(tcell.ColorBlack)
	}
}
