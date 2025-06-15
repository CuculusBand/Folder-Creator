package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// createFolderStructure()
	MyApp := app.New()
	// Load and Set the custom font file
	customFont := fyne.NewStaticResource("NotoSans", LoadFont("assets/font/NotoSans-SemiBold.ttf"))
	MyApp.Settings().SetTheme(&appTheme{font: customFont})
	mainWindow := MyApp.NewWindow("Folder Creator")
	mainWindow.Resize(fyne.NewSize(1280, 800))
	mainWindow.SetFixedSize(true)
	//mainWindow.SetContent(container.NewVBox(widget.NewLabel("<Folder Creator>\n")))
	FileOpenDialog := widget.NewButton("Open File Dialog", func() {
		log.Println("Open File Dialog button clicked")
		// Here you would implement the logic to open a file dialog
		// and handle the selected file path.
	})
	content2 := widget.NewButton("destinationPath", func() {
		log.Println("destinationPathPath button clicked")
	})
	// Add both buttons to the main window using a container
	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("<Folder Creator>\n"),
		FileOpenDialog,
		content2,
	))
	// Run the application
	mainWindow.Show()
	print("Folder Creator is running...\n")
	MyApp.Run()
	print("Folder Creator is closed.\n")
}
