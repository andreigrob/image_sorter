package main

import (
	"log"
	"os"

	ft "fmt"
	fp "path/filepath"
)

type Model struct {
	c         *Controller
	ImagesDir FileName
	images    []FileName
	index     int
}

func (m *Model) SetController(c *Controller) {
	m.c = c
}

func (m *Model) Init(imagesDir string) {
	m.ImagesDir = FileName(imagesDir)
	m.index = 0
	files, err := os.ReadDir(imagesDir)
	if err != nil {
		log.Fatalf("Error reading directory: %v\n", err)
	}

	// filter the files to only include pngs and jpgs
	imageFiles := make([]FileName, 0, len(files))
	var name FileName
	var i int
	for ; i < len(files); i++ {
		if files[i].IsDir() {
			continue
		}
		if name = FileName(files[i].Name()); name.HasExtension(`.png`, `.jpg`, `.jpeg`) {
			imageFiles = append(imageFiles, FileName(fp.Join(imagesDir, string(name))))
		}
	}
	m.images = imageFiles
}

func (m *Model) CurrentImage() (_ FileName) {
	return m.images[m.index]
}

func (m *Model) GetImageNum() (_ int) {
	return m.index + 1
}

func (m *Model) GetLenImages() (_ int) {
	return len(m.images)
}

func (m *Model) GetName() (_ string) {
	return "Image Sorter"
}

func (m *Model) GetTitle() (_ string) {
	return ft.Sprintf("Image Viewer [%d/%d]", m.GetImageNum(), len(m.images))
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
