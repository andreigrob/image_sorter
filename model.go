package main

import (
	ft "fmt"
	"log"
)

type DriveImage struct {
	ID   string
	Name string
	Path string
}

type Model struct {
	c        *Controller
	SourceID string
	images   []DriveImage
	index    int
}

func (m *Model) SetController(c *Controller) {
	m.c = c
}

func (m *Model) Init(folderID string) {
	m.SourceID = folderID
	m.index = 0

	files, err := m.c.D.ListImages(folderID)
	if err != nil {
		log.Fatalf("error listing drive folder: %v", err)
	}

	m.images = make([]DriveImage, len(files))
	for i, f := range files {
		m.images[i] = DriveImage{ID: f.Id, Name: f.Name}
	}
}

func (m *Model) CurrentImage() (_ FileName) {
	if len(m.images) == 0 {
		return ""
	}
	img := &m.images[m.index]
	if img.Path == "" {
		path, err := m.c.D.DownloadFile(img.ID, img.Name)
		if err != nil {
			log.Printf("error downloading file: %v", err)
			return ""
		}
		img.Path = path
	}
	return FileName(img.Path)
}

func (m *Model) GetImageNum() (_ int) {
	return m.index + 1
}

func (m *Model) GetLenImages() (_ int) {
	return len(m.images)
}

func (m *Model) GetName() (_ string) {
	return `Image Sorter`
}

func (m *Model) GetTitle() (_ string) {
	if len(m.images) == 0 {
		return m.GetName() + " [0/0]"
	}
	return ft.Sprintf("%s [%d/%d]", m.GetName(), m.GetImageNum(), len(m.images))
}

func (m *Model) Next() (_ bool) {
	if m.index < len(m.images)-1 {
		m.index++
		return true
	}
	return
}

func (m *Model) Prev() (_ bool) {
	if m.index > 0 {
		m.index--
		return true
	}
	return
}

func (m *Model) MoveCurrentImage(folder string) (_ bool) {
	if len(m.images) == 0 {
		return
	}

	destID, ok := m.c.D.Destinations[folder]
	if !ok {
		log.Printf("destination folder %s not configured", folder)
		return
	}

	img := m.images[m.index]
	if err := m.c.D.MoveFile(img.ID, destID); err != nil {
		log.Printf("error moving file: %v", err)
		return
	}

	m.images = append(m.images[:m.index], m.images[m.index+1:]...)
	if m.index >= len(m.images) {
		m.index--
	}
	return true
}
