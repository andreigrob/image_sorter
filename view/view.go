package view

import (
	st "strings"

	fy "fyne.io/fyne/v2"
	cv "fyne.io/fyne/v2/canvas"
	cn "fyne.io/fyne/v2/container"
	wg "fyne.io/fyne/v2/widget"

	in "github.com/andreigrob/image_sorter/interfaces"
	tp "github.com/andreigrob/image_sorter/types"
)

type DataViewT struct {
	c          in.IController
	window     fy.Window
	img        *cv.Image
	titleLabel *wg.Label
	imageLabel *wg.Label
}

type ViewT struct {
	*DataViewT
}

func New(c in.IController, a fy.App) (v ViewT) {
	v.DataViewT = &DataViewT{}
	v.SetController(c)

	v.titleLabel = wg.NewLabel(v.c.GetName())
	v.imageLabel = wg.NewLabel(v.c.GetTitle())

	currentImage, _ := v.c.GetCurrentImage()
	v.img = cv.NewImageFromFile(string(currentImage))
	v.img.FillMode = cv.ImageFillContain

	labels := cn.NewVBox(
		v.titleLabel,
		v.imageLabel,
	)
	navButtons := cn.NewHBox(
		wg.NewButton(`Prev`, v.c.ShowPrevImage),
		wg.NewButton(`Next`, v.c.ShowNextImage),
	)
	sortButtons := cn.NewHBox(v.MoveButtons(buttonNames...)...)
	buttons := cn.NewVBox(navButtons, sortButtons)
	content := cn.NewBorder(labels, buttons, nil, nil, v.img)

	v.window = a.NewWindow(``)
	v.window.SetContent(content)
	v.window.Resize(fy.NewSize(800, 600))

	return
}

func (v ViewT) SetController(c in.IController) {
	v.c = c
}

func (v ViewT) MoveButton(folder string) (_ *wg.Button) {
	folderName := st.ToLower(folder)
	return wg.NewButton(folder,
		func() {
			_ = v.c.MoveImage(tp.NameT(folderName))
		},
	)
}

func (v ViewT) MoveButtons(folders ...string) (_ []fy.CanvasObject) {
	buttons := make([]fy.CanvasObject, len(folders))
	var i int = len(folders) - 1
	for ; i >= 0; i-- {
		buttons[i] = v.MoveButton(folders[i])
	}
	return buttons
}

var buttonNames = []string{`Glass`, `Metal`, `Paper`, `Plastic`}

func (v ViewT) SetTitle(title string) {
	v.window.SetTitle(title)
	v.imageLabel.SetText(title)
}

func (v ViewT) ShowImage() {
	v.SetTitle(v.c.GetTitle())
	currentName, _ := v.c.GetCurrentImage()
	currentImage := string(currentName)
	if currentImage == `` {
		v.img.File = `no_images.png`
	} else {
		v.img.File = currentImage
	}
	v.img.Refresh()
}

func (v ViewT) Display() {
	v.ShowImage()
	v.window.ShowAndRun()
}
