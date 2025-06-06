package main

import (
	"log"
	"os"

	fa "fyne.io/fyne/v2/app"
)

func main() {
	// read the images directory from the command line
	if len(os.Args) < 2 {
		log.Fatal("images directory is required")
	}

	imageViewer := &Controller{}
	imageViewer.Init()

	imageViewer.M.Init(os.Args[1])

	myApp := fa.New()
	imageViewer.V.Init(myApp)

	imageViewer.Start()
}
