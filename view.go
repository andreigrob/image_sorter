package main

import (
	fy "fyne.io/fyne/v2"
	fcv "fyne.io/fyne/v2/canvas"
	fct "fyne.io/fyne/v2/container"
	fwg "fyne.io/fyne/v2/widget"
)

type View struct {
	c      *Controller
	window fy.Window
	img    *fcv.Image
	label  *fwg.Label
	label2 *fwg.Label
}

func (v *View) SetController(c *Controller) {
	v.c = c
}

func (v *View) Init(a fy.App) {
	v.label = fwg.NewLabel(v.c.M.GetName())
	v.label2 = fwg.NewLabel(v.c.M.GetTitle())

	v.img = fcv.NewImageFromFile(string(v.c.M.CurrentImage()))
	v.img.FillMode = fcv.ImageFillContain

	content := fct.NewVBox(
		v.label,
		v.label2,
		v.img,
		fct.NewHBox(
			fwg.NewButton(`Prev`, v.c.ShowPrevImage),
			fwg.NewButton(`Next`, v.c.ShowNextImage),
		),
	)

	v.window = a.NewWindow(string(v.c.M.ImagesDir))
	v.window.SetContent(content)
	v.window.Resize(fy.NewSize(800, 600))
}

func (v *View) SetTitle(title string) {
	v.window.SetTitle(title)
	v.label2.SetText(title)
}

func (v *View) ShowImage() {
	v.SetTitle(v.c.M.GetTitle())
	v.img.File = string(v.c.M.CurrentImage())
	v.img.Refresh()
}

func (v *View) Display() {
	v.ShowImage()
	v.window.ShowAndRun()
}
