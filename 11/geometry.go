package main

type coord struct {
	x int
	y int
}

type rect struct {
	topLeft coord
	bottomRight coord
}

func (coord) boundingBox(c1 coord, c2 coord) rect {
	minX := Min(c1.x, c2.x)
	minY := Min(c1.y, c2.y)
	maxX := Max(c1.x, c2.x)
	maxY := Max(c1.y, c2.y)

	return rect{topLeft: coord{
		x: minX,
		y: minY,
	},
		bottomRight: coord{
			x: maxX,
			y: maxY,
		}}
}

type compassDirection struct {
	xPart int
	yPart int
}

func (c *compassDirection) rotateAnticlockwise() {
	c.xPart, c.yPart = c.yPart, -c.xPart
}

func (c *compassDirection) rotateClockwise() {
	c.xPart, c.yPart = -c.yPart, c.xPart
}
