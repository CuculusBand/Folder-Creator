package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// MainApp holds the main application structure, including the app, window, and file processor
type MainApp struct {
	App          fyne.App
	Window       fyne.Window
	Processor    *FileProcessor
	PreviewTable *widget.Table
	StatusLabel  *widget.Label
	FileLabel    *widget.Label
	DestLabel    *widget.Label
}

// InitializeApp holds the application and window instances along with a file processor
func InitializeApp(app fyne.App, window fyne.Window) *MainApp {
	return &MainApp{
		App:       app,
		Window:    window,
		Processor: NewFileProcessor(),
	}
}

func (a *MainApp) MakeUI() {
	// Create buttons
	FileSelectButton := widget.NewButton("Select File", a.selectTableFile)
	TargetSelectButton := widget.NewButton("Target Path", a.selectDestination)
	ClearButton := widget.NewButton("Clear", a.clearAll)
	CreateButton := widget.NewButton("Create", a.generateFolders)
	ExitButton := widget.NewButton("Exit", func() { a.App.Quit() })
	// Button layout
	buttonRow := container.NewHBox(
		FileSelectButton,
		TargetSelectButton,
		layout.NewSpacer(),
		ClearButton,
		CreateButton,
		ExitButton,
	)
	// Create Lables
	a.StatusLabel = widget.NewLabel("Ready")
	a.StatusLabel.Wrapping = fyne.TextWrapWord
	a.FileLabel = widget.NewLabel("No file selected")
	a.DestLabel = widget.NewLabel("No destination selected")
	// File information area
	fileInfo := container.NewGridWithColumns(2,
		widget.NewLabel("Table file:"), a.FileLabel,
		widget.NewLabel("Targer path:"), a.DestLabel,
	)
	// Create preview table
	a.PreviewTable = widget.NewTable(
		func() (int, int) {
			if len(a.Processor.TableData) == 0 {
				return 1, 1
			}
			return len(a.Processor.TableData), len(a.Processor.TableData[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if len(a.Processor.TableData) > i.Row && len(a.Processor.TableData[i.Row]) > i.Col {
				label.SetText(a.Processor.TableData[i.Row][i.Col])
			} else {
				label.SetText("")
			}
		},
	)
	// Config table's cloumn widths
	a.PreviewTable.SetColumnWidth(0, 170)
	for i := 1; i < 20; i++ {
		a.PreviewTable.SetColumnWidth(i, 170)
	}
	// Create the main content layout
	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("|Folder Creator|"),
			widget.NewSeparator(),
			fileInfo,
			widget.NewSeparator(),
			buttonRow,
			widget.NewSeparator(),
			widget.NewLabel("Preview:"),
		),
		a.StatusLabel,
		nil,
		nil,
		container.NewScroll(a.PreviewTable),
	)
	// Set the content
	a.Window.SetContent(content)
}

// Select a file to load table data
func (a *MainApp) selectTableFile() {
	dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			a.StatusLabel.SetText("Wrong file: " + err.Error())
			return
		}
		if reader == nil {
			return
		}
		FilePath := reader.URI().Path()
		if runtime.GOOS == "windows" {
			// Remove leading slash for Windows paths
			if len(FilePath) > 2 && FilePath[0] == '/' && FilePath[2] == ':' {
				FilePath = FilePath[1:]
			}
			// Replace forward slashes with backslashes for Windows compatibility
			FilePath = strings.ReplaceAll(FilePath, "/", "\\")
		}

		a.FileLabel.SetText(filepath.Base(FilePath))
		a.StatusLabel.SetText("Loading...")
		// Load the file
		if err := a.Processor.LoadFile(FilePath); err != nil {
			a.StatusLabel.SetText("Failed to load: " + err.Error())
			return
		}
		a.PreviewTable.Refresh()
		a.StatusLabel.SetText(fmt.Sprintf("All data loaded: %d rows", len(a.Processor.TableData)))
	}, a.Window).Show()
}

// Select a destination folder to create new folders
func (a *MainApp) selectDestination() {
	dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
		if err != nil {
			a.StatusLabel.SetText("Wrong target path: " + err.Error())
			return
		}
		if list == nil {
			return
		}
		//RawDestPath := list.Path()
		a.Processor.DestPath = list.Path()
		a.DestLabel.SetText(a.Processor.DestPath)
		a.StatusLabel.SetText("Selected target path: " + filepath.Base(a.Processor.DestPath))
	}, a.Window).Show()
}

// Clear all content in the table
func (a *MainApp) clearAll() {
	a.Processor.Clear()
	a.FileLabel.SetText("No file selected")
	a.DestLabel.SetText("No destination selected")
	a.StatusLabel.SetText("All content cleared")
	a.PreviewTable.Refresh()
}

// Create folders based on the loaded table data
func (a *MainApp) generateFolders() {
	if a.Processor.TableFilePath == "" {
		a.StatusLabel.SetText("Select a file first!")
		return
	}
	if a.Processor.DestPath == "" {
		a.StatusLabel.SetText("Select a target path first!")
		return
	}
	if len(a.Processor.TableData) == 0 {
		a.StatusLabel.SetText("No available data!")
		return
	}

	successCount, err := a.Processor.GenerateFolders()
	if err != nil {
		a.StatusLabel.SetText("Error: " + err.Error())
		return
	}
	a.StatusLabel.SetText(fmt.Sprintf("Sucessfully created %d folder(s)", successCount))
}
