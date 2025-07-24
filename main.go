package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Create the application
	MyApp := app.NewWithID("Folder Creator")
	// Load and Set the custom font file
	//customFont := fyne.NewStaticResource("NotoSans", LoadFont("fonts/NotoSans-SemiBold.ttf"))
	MyApp.Settings().SetTheme(&appTheme{regularFont: AppFont})
	// Set the application icon from a byte slice
	var iconData []byte
	icon := fyne.NewStaticResource("GO-FolderCreator-icon.png", iconData)
	MyApp.SetIcon(icon)

	// Create the window
	MainWindow := MyApp.NewWindow("Folder Creator")
	MainWindow.Resize(fyne.NewSize(600, 850))
	MainWindow.SetFixedSize(false)
	MainWindow.CenterOnScreen()
	app := InitializeApp(MyApp, MainWindow)
	time.Sleep(50 * time.Millisecond)
	// Create UI
	app.MakeUI()
	// Run the application
	MainWindow.ShowAndRun()
}
