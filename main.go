package main

import (
	ct "context"
	"log"
	"os"

	fa "fyne.io/fyne/v2/app"
)

func main() {
	// Expect credentials file and source folder ID
	if len(os.Args) < 3 {
		log.Fatal("usage: image_sorter <credentials.json> <sourceFolderName>")
	}

	credFile := os.Args[1]
	folderName := os.Args[2]

	ctx := ct.Background()
	imageViewer := &Controller{}
	imageViewer.Init()

	e := imageViewer.D.Init(ctx, credFile)
	if e != nil {
		log.Fatalf("unable to create drive service: %v", e)
	}

	folderID, e := imageViewer.D.FindFolderID(folderName)
	if e != nil {
		log.Fatalf("unable to find folder %s: %v", folderName, e)
	}
	imageViewer.M.Init(folderID)

	// initialize known destination folders if they exist
	folders := []string{`glass`, `metal`, `paper`, `plastic`}
	var id string
	var i int
	for ; i < len(folders); i++ {
		id, e = imageViewer.D.FindFolderID(folders[i])
		if e != nil {
			log.Printf("warning: folder %s not found: %v", folders[i], e)
			continue
		}
		imageViewer.D.Destinations[folders[i]] = id
	}

	myApp := fa.New()
	imageViewer.V.Init(myApp)

	imageViewer.Start()
}
