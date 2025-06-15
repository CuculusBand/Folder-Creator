package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Create the application
	MyApp := app.NewWithID("FolderCreator")
	// Load and Set the custom font file
	customFont := fyne.NewStaticResource("NotoSans", LoadFont("assets/font/NotoSans-SemiBold.ttf"))
	MyApp.Settings().SetTheme(&appTheme{font: customFont})
	// Create the window
	MainWindow := MyApp.NewWindow("Folder Creator")
	MainWindow.Resize(fyne.NewSize(600, 800))
	MainWindow.SetFixedSize(true)
	MainWindow.CenterOnScreen()
	app := InitializeApp(MyApp, MainWindow)
	// Create UI
	app.MakeUI()
	print("Folder Creator is running...\n")
	// Run the application
	MainWindow.ShowAndRun()
	print("Folder Creator is closed\n")
}
