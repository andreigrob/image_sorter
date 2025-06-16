package main

import (
	ct "context"
	"fmt"
	"log"
	"os"

	fa "fyne.io/fyne/v2/app"

	cr "github.com/andreigrob/image_sorter/controller"
	ml "github.com/andreigrob/image_sorter/model"
	vw "github.com/andreigrob/image_sorter/view"

	gd "github.com/andreigrob/image_sorter/gdrive"

	tp "github.com/andreigrob/image_sorter/types"
)

var folders = []tp.NameT{`glass`, `metal`, `paper`, `plastic`}

func Init(ctx ct.Context, imageViewer *cr.ControllerT, credFile string, folderName string) (e error) {
	// initialize GDrive service
	log.Println("Initializing GDrive service")
	gDrive, e := gd.New(ctx, credFile)
	if e != nil {
		fmt.Printf("Error creating GDrive service: %v\n", e)
		return
	}
	imageViewer.D = gDrive

	// initialize the model
	log.Println("Initializing model")
	model, e := ml.New(imageViewer, folderName)
	if e != nil {
		fmt.Printf("Error initializing model: %v\n", e)
		return
	}
	imageViewer.M = model

	// find destination folders on GDrive
	log.Println("Finding destination folders on GDrive")
	_ = imageViewer.FindFolders(folders)

	return
}

func main() {
	// Expect credentials file and source folder ID
	if len(os.Args) < 3 {
		log.Println("usage: image_sorter <credentials.json> <sourceFolderName>")
		return
	}

	imageViewer := cr.New()

	ctx, cancel := ct.WithCancel(ct.Background())
	defer cancel()

	credFile := os.Args[1]
	folderName := os.Args[2]
	e := Init(ctx, &imageViewer, credFile, folderName)
	if e != nil {
		log.Printf("Error initializing image viewer: %v\n", e)
		return
	}

	// initialize view
	viewApp := fa.New()
	imageViewer.V = vw.New(&imageViewer, viewApp)

	// start the application
	imageViewer.Start()
}
