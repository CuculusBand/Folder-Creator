package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// createFolderStructure()
	app := app.New()
	w0 := app.NewWindow("Folder Creator")
	w0.Resize(fyne.NewSize(1280, 800))
	w0.SetFixedSize(true)
	w0.SetContent(container.NewVBox(widget.NewLabel("<Folder Creator>\n")))
	w0.Show()
	print("Folder Creator is running...\n")
	app.Run()
	print("Folder Creator is closed.\n")
}
