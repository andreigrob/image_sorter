package interfaces

import (
	gd "github.com/andreigrob/image_sorter/gdrive"
	tp "github.com/andreigrob/image_sorter/types"
	ut "github.com/andreigrob/image_sorter/utils"
)

type IController interface {
	GetDrive() (_ gd.GDriveT)

	ShowNextImage()
	ShowPrevImage()
	MoveImage(folder tp.NameT) (_ error)

	GetName() (_ string)
	GetTitle() (_ string)
	GetCurrentImage() (_ ut.FileName, _ error)
}
