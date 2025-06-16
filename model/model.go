package model

import (
	"errors"
	"fmt"
	"log"

	gd "github.com/andreigrob/image_sorter/gdrive"
	in "github.com/andreigrob/image_sorter/interfaces"
	tp "github.com/andreigrob/image_sorter/types"
	ut "github.com/andreigrob/image_sorter/utils"
)

type DataModelT struct {
	cr     in.IController
	gDrive gd.GDriveT
	folder tp.GFileT
	images []tp.LocalGFileT
	index  int
}

type ModelT struct {
	*DataModelT
}

func (m ModelT) SetController(c in.IController) {
	m.cr = c
	m.gDrive = m.cr.GetDrive()
}

func New(c in.IController, folderName string) (m ModelT, e error) {
	m.DataModelT = &DataModelT{}

	m.SetController(c)

	m.folder.Name = tp.NameT(folderName)
	m.folder.Id, e = m.gDrive.FindFolderId(m.folder.Name)
	if e != nil {
		fmt.Printf("Folder %s not found on GDrive: %v\n", m.folder.Name, e)
		return
	}

	fmt.Printf("Folder %v found on GDrive\n", m.folder)

	m.index = 0

	files, e := m.gDrive.ListImages(m.folder)
	if e != nil {
		log.Printf("Error listing GDrive folder %v: %v", m.folder, e)
		return
	}

	m.images = make([]tp.LocalGFileT, len(files))
	var i int = len(files) - 1
	for ; i >= 0; i-- {
		m.images[i].Id = tp.IdT(files[i].Id)
		m.images[i].Name = tp.NameT(files[i].Name)
	}

	return
}

func (m ModelT) CurrentImage() (_ ut.FileName, e error) {
	if len(m.images) == 0 {
		return
	}
	img := &m.images[m.index]
	if img.Path == "" {
		img.Path, e = m.gDrive.Download(img.GFileT)
		if e != nil {
			log.Printf("Error downloading %s (%s): %v", img.Name, img.Id, e)
			return
		}
	}

	return ut.FileName(img.Path), e
}

func (m ModelT) GetImageNum() (_ int) {
	return m.index + 1
}

func (m ModelT) GetImageLen() (_ int) {
	return len(m.images)
}

func (m ModelT) GetName() (_ string) {
	return `Image Sorter`
}

func (m ModelT) GetTitle() (_ string) {
	if len(m.images) == 0 {
		return m.GetName() + " [0/0]"
	}
	return fmt.Sprintf("%s [%d/%d]", m.GetName(), m.GetImageNum(), m.GetImageLen())
}

func (m ModelT) Next() (_ bool) {
	if len(m.images)-m.index > 1 {
		m.index++
		return true
	}
	return
}

func (m ModelT) Prev() (_ bool) {
	if m.index > 0 {
		m.index--
		return true
	}
	return
}

var ErrNoDestinationFolder = errors.New("destination folder not configured")

func (m ModelT) MoveCurrentImage(folder tp.NameT) (ok bool, e error) {
	if len(m.images) == 0 {
		return
	}

	dest := tp.GFileT{Name: folder}
	dest.Id, ok = m.gDrive.Ids[dest.Name]
	if !ok {
		log.Printf("Destination folder %s not configured", dest.Name)
		return ok, ErrNoDestinationFolder
	}

	img := m.images[m.index]
	if e = m.gDrive.Move(img.GFileT, dest); e != nil {
		log.Printf("Error moving file %v to folder %v: %v", img, dest, e)
		return
	}

	m.images = append(m.images[:m.index], m.images[m.index+1:]...)
	if m.index >= len(m.images) {
		m.index--
	}

	return true, e
}
