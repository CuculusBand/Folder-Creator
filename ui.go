package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MainApp holds the main application structure, including the app, window, and file processor
type MainApp struct {
	App          fyne.App
	Window       fyne.Window
	Processor    *FileProcessor
	PreviewTable *widget.Table
	StatusLabel  *widget.Label
	FilePath     *PathDisplay
	DestPath     *PathDisplay
	ThemeButton  *widget.Button
	DarkMode     bool
}

// PathDisplay shows the file or folder path in a scrollable text container
type PathDisplay struct {
	Text      *canvas.Text
	Container *container.Scroll
}

// InitializeApp holds the application and window instances along with a file processor
func InitializeApp(app fyne.App, window fyne.Window) *MainApp {
	return &MainApp{
		App:       app,
		Window:    window,
		Processor: NewFileProcessor(),
		DarkMode:  false,
	}
}

// Sets up the UI for the application
func (a *MainApp) MakeUI() {
	// Add theme control button, refreshes the theme when clicked
	a.ThemeButton = widget.NewButton("🌙", a.ToggleTheme)
	// Set the theme when the app starts
	a.SetTheme(a.DarkMode)
	// Create buttons
	FileSelectButton := widget.NewButton("Select File", a.SelectTableFile)
	TargetSelectButton := widget.NewButton("Target Path", a.SelectDestination)
	ClearButton := widget.NewButton("Clear", a.ClearAll)
	CreateButton := widget.NewButton("Create", a.GenerateFolders)
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
	// Set the title of the app
	Title := widget.NewLabel("<Folder Creator>")
	// Title and theme button layout
	TitleContainer := container.NewHBox(
		Title,
		layout.NewSpacer(),
		a.ThemeButton,
	)
	// Create scrollable path displays
	a.FilePath = CreatePathDisplay()
	a.DestPath = CreatePathDisplay()
	// Create status Lables
	a.StatusLabel = widget.NewLabel("Ready")
	a.StatusLabel.Wrapping = fyne.TextWrapWord
	// File information area
	fileInfo := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Table file:	"),
			a.FilePath.Container,
		),
		container.NewHBox(
			widget.NewLabel("Target path:	"),
			a.DestPath.Container,
		),
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
	a.PreviewTable.SetColumnWidth(0, 150)
	for i := 1; i < 20; i++ {
		a.PreviewTable.SetColumnWidth(i, 150)
	}
	// Create the main content layout
	content := container.NewBorder(
		container.NewVBox(
			TitleContainer,
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

// Use canvas to display file paths
func CreatePathDisplay() *PathDisplay {
	// Set text first
	text := canvas.NewText("No Selection", color.Black)
	text.TextSize = 14
	text.TextStyle = fyne.TextStyle{Monospace: false, Bold: true}
	// Create a scrollable container for the text
	scroll := container.NewHScroll(text)
	// Set min size for labels and add scrolls
	scroll.SetMinSize(fyne.NewSize(475, 45))
	return &PathDisplay{
		Text:      text,
		Container: scroll,
	}
}

// Refreshes PathDisplay's text color based on the theme
func (pd *PathDisplay) RefreshColor(isDark bool) {
	if isDark {
		pd.Text.Color = color.White // Use White for dark theme
	} else {
		pd.Text.Color = color.Black // Use Black for light theme
	}
	pd.Text.Refresh()
}

// Select a file to load table data
func (a *MainApp) SelectTableFile() {
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
		// Set the file path to the label
		a.FilePath.Text.Text = FilePath
		a.FilePath.Text.Refresh()
		// a.FilePath.Path.Text = filepath.Base(FilePath)
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
func (a *MainApp) SelectDestination() {
	dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
		if err != nil {
			a.StatusLabel.SetText("Wrong target path: " + err.Error())
			return
		}
		if list == nil {
			return
		}
		a.Processor.DestPath = list.Path()
		a.DestPath.Text.Text = a.Processor.DestPath
		a.DestPath.Text.Refresh()
		a.StatusLabel.SetText("Selected target path: " + filepath.Base(a.Processor.DestPath))
	}, a.Window).Show()
}

// Clear all content in the table
func (a *MainApp) ClearAll() {
	a.Processor.Clear()
	a.FilePath.Text.Text = "No Selection"
	a.FilePath.Text.Refresh()
	a.DestPath.Text.Text = "No Selection"
	a.DestPath.Text.Refresh()
	a.StatusLabel.SetText("All content cleared")
	a.PreviewTable.Refresh()
}

// Create folders based on the loaded table data
func (a *MainApp) GenerateFolders() {
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

// Toggle the theme between light and dark mode
func (a *MainApp) SetTheme(isDark bool) {
	a.DarkMode = isDark
	// Save the theme preference
	// Set Theme based on the button state
	a.App.Preferences().SetBool("dark_mode", isDark)
	if isDark {
		a.App.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.App.Settings().SetTheme(theme.LightTheme())
	}
}

// Toggle the theme when the button is clicked
func (a *MainApp) ToggleTheme() {
	a.SetTheme(!a.DarkMode)
	if a.DarkMode {
		a.ThemeButton.SetText("☀️") // Show sun icon if dark mode is enabled
	} else {
		a.ThemeButton.SetText("🌙") // Show moon icon if dark mode is disabled
	}
	a.Window.Content().Refresh()
	// Update PathDisplays's colors
	a.FilePath.RefreshColor(a.DarkMode)
	a.DestPath.RefreshColor(a.DarkMode)
}
