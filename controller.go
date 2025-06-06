package main

type Controller struct {
	M Model
	V View
}

func (c *Controller) Start() {
	c.V.Display()
}

func (c *Controller) Init() {
	c.M.SetController(c)
	c.V.SetController(c)
}

func (c *Controller) ShowNextImage() {
	if c.M.Next() {
		c.V.ShowImage()
	}
}

func (c *Controller) ShowPrevImage() {
	if c.M.Prev() {
		c.V.ShowImage()
	}
}

func (c *Controller) MoveImage(folder string) {
	if c.M.MoveCurrentImage(folder) {
		c.V.ShowImage()
	}
}
