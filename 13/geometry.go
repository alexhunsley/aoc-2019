package main

import "math"

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

type boundsTracker struct {
	minX int
	minY int
	maxX int
	maxY int
}

func (b *boundsTracker) Init() {
	b.maxX = math.MinInt64
	b.maxY = math.MinInt64
	b.minX = math.MaxInt64
	b.minY = math.MaxInt64
}

func (b *boundsTracker) addCoord(x int, y int) {
	if x < b.minX {
		b.minX = x
	}
	if x > b.maxX {
		b.maxX = x
	}
	if y < b.minY {
		b.minY = y
	}
	if y > b.maxY {
		b.maxY = y
	}
}
