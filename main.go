package main

import (
	"image/color"
	"log"
	"os"

	fy "fyne.io/fyne/v2"
	fa "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

type largeTextTheme struct{}

func (m largeTextTheme) Color(n fy.ThemeColorName, v fy.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (m largeTextTheme) Font(s fy.TextStyle) fy.Resource {
	return theme.DefaultTheme().Font(s)
}

func (m largeTextTheme) Icon(n fy.ThemeIconName) fy.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (m largeTextTheme) Size(n fy.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return 38 // Adjust this number for larger or smaller text
	}
	return theme.DefaultTheme().Size(n)
}

func main() {
	// read the images directory from the command line
	if len(os.Args) < 2 {
		log.Fatal("images directory is required")
	}

	imageViewer := &Controller{}
	imageViewer.Init()

	imageViewer.M.Init(os.Args[1])
	myApp := fa.NewWithID("image-viewer")
	myApp.Settings().SetTheme(&largeTextTheme{}) // Apply custom theme

	imageViewer.V.Init(myApp)

	imageViewer.Start()
}
