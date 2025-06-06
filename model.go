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
	if len(m.images) == 0 {
		return ""
	}
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

	// create the folder if it doesn't exist
	e := os.MkdirAll(folder, 0o755)
	if e != nil {
		log.Printf("Error creating folder %s: %v\n", folder, e)
		return
	}

	current := string(m.CurrentImage())
	dest := fp.Join(folder, fp.Base(current))
	if e = os.Rename(current, dest); e != nil {
		log.Printf("Error moving file: %v\n", e)
		return
	}

	m.images = append(m.images[:m.index], m.images[m.index+1:]...)
	if m.index >= len(m.images) {
		m.index--
	}
	return true
}
