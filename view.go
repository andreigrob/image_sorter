package main

import (
	st "strings"

	fy "fyne.io/fyne/v2"
	fcv "fyne.io/fyne/v2/canvas"
	fct "fyne.io/fyne/v2/container"
	fwg "fyne.io/fyne/v2/widget"
)

type View struct {
	c          *Controller
	window     fy.Window
	img        *fcv.Image
	titleLabel *fwg.Label
	imageLabel *fwg.Label
}

func (v *View) SetController(c *Controller) {
	v.c = c
}

func (v *View) MoveButton(folder string) (_ *fwg.Button) {
	folderName := st.ToLower(folder)
	return fwg.NewButton(folder, func() { v.c.MoveImage(folderName) })
}

func (v *View) MoveButtons(folders ...string) (_ []fy.CanvasObject) {
	buttons := make([]fy.CanvasObject, len(folders))
	var i int
	for ; i < len(folders); i++ {
		buttons[i] = v.MoveButton(folders[i])
	}
	return buttons
}

func (v *View) Init(a fy.App) {
	v.titleLabel = fwg.NewLabel(v.c.M.GetName())
	v.imageLabel = fwg.NewLabel(v.c.M.GetTitle())

	v.img = fcv.NewImageFromFile(string(v.c.M.CurrentImage()))
	v.img.FillMode = fcv.ImageFillContain

	labels := fct.NewVBox(
		v.titleLabel,
		v.imageLabel,
	)
	navButtons := fct.NewHBox(
		fwg.NewButton(`Prev`, v.c.ShowPrevImage),
		fwg.NewButton(`Next`, v.c.ShowNextImage),
	)
	sortButtons := fct.NewHBox(v.MoveButtons(`Glass`, `Metal`, `Paper`, `Plastic`)...)
	buttons := fct.NewVBox(navButtons, sortButtons)
	content := fct.NewBorder(labels, buttons, nil, nil, v.img)

	v.window = a.NewWindow(``)
	v.window.SetContent(content)
	v.window.Resize(fy.NewSize(800, 600))
}

func (v *View) SetTitle(title string) {
	v.window.SetTitle(title)
	v.imageLabel.SetText(title)
}

func (v *View) ShowImage() {
	v.SetTitle(v.c.M.GetTitle())
	currentImage := string(v.c.M.CurrentImage())
	if currentImage == `` {
		v.img.File = `no_images.png`
	} else {
		v.img.File = currentImage
	}
	v.img.Refresh()
}

func (v *View) Display() {
	v.ShowImage()
	v.window.ShowAndRun()
}
