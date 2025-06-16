package controller

import (
	"fmt"

	gd "github.com/andreigrob/image_sorter/gdrive"

	tp "github.com/andreigrob/image_sorter/types"
	ut "github.com/andreigrob/image_sorter/utils"

	ml "github.com/andreigrob/image_sorter/model"
	vw "github.com/andreigrob/image_sorter/view"
)

type DataControllerT struct {
	M ml.ModelT
	V vw.ViewT
	D gd.GDriveT
}

type ControllerT struct {
	*DataControllerT
}

func New() (c ControllerT) {
	c.DataControllerT = &DataControllerT{}
	return
}

func (c ControllerT) Start() {
	c.V.Display()
}

func (c ControllerT) GetDrive() (_ gd.GDriveT) {
	return c.D
}

func (c ControllerT) ShowNextImage() {
	if c.M.Next() {
		c.V.ShowImage()
	}
}

func (c ControllerT) ShowPrevImage() {
	if c.M.Prev() {
		c.V.ShowImage()
	}
}

func (c ControllerT) MoveImage(folder tp.NameT) (e error) {
	ok, e := c.M.MoveCurrentImage(folder)
	if e != nil {
		currentImage, _ := c.GetCurrentImage()
		fmt.Printf("Error moving %s to %s: %v\n", currentImage, folder, e)
		return
	}
	if ok {
		c.V.ShowImage()
	}
	return
}

func (c ControllerT) GetName() string {
	return c.M.GetName()
}

func (c ControllerT) GetTitle() string {
	return c.M.GetTitle()
}

func (c ControllerT) GetCurrentImage() (_ ut.FileName, e error) {
	return c.M.CurrentImage()
}

// FindFolders finds the ids of folders on the GDrive and stores them in the Ids map.
func (c ControllerT) FindFolders(folders []tp.NameT) (e error) {
	var folderId tp.IdT
	var i int
	for ; i < len(folders); i++ {
		folderId, e = c.D.FindFolderId(folders[i])
		if e != nil {
			fmt.Printf("Folder %s not found on GDrive: %v\n", folders[i], e)
			continue
		}

		c.D.Ids[folders[i]] = folderId
	}
	return
}
