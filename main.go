package main

import (
        "context"
        "log"
        "os"

        fa "fyne.io/fyne/v2/app"
)

func main() {
        // Expect credentials file and source folder ID
        if len(os.Args) < 3 {
                log.Fatal("usage: image_sorter <credentials.json> <sourceFolderID>")
        }

        credFile := os.Args[1]
        folderID := os.Args[2]

        ctx := context.Background()
        driveSvc, err := NewDrive(ctx, credFile)
        if err != nil {
                log.Fatalf("unable to create drive service: %v", err)
        }

        // initialize known destination folders if they exist
        for _, name := range []string{"glass", "metal", "paper", "plastic"} {
                id, err := driveSvc.FindFolderID(name)
                if err != nil {
                        log.Printf("warning: folder %s not found: %v", name, err)
                        continue
                }
                driveSvc.Destinations[name] = id
        }

        imageViewer := &Controller{D: driveSvc}
        imageViewer.Init()

        imageViewer.M.Init(folderID)

	myApp := fa.New()
	imageViewer.V.Init(myApp)

	imageViewer.Start()
}
