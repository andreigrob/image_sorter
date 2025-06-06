package main

import (
	fa "fyne.io/fyne/v2/app"
)

func main() {
	imageViewer := &Controller{}
	imageViewer.Init()

	const imagesDir = `images/`
	imageViewer.M.Init(imagesDir)

	myApp := fa.New()
	imageViewer.V.Init(myApp)

	imageViewer.Start()
}
