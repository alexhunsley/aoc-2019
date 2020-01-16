package main

type paintBot struct {
	// the colour painted at coords
	colour map[coord]int
	facingDirection compassDirection
	location coord
}

func (p *paintBot) moveForward() {
	p.location.x += p.facingDirection.xPart
	p.location.y += p.facingDirection.yPart
}

func (p *paintBot) getColour() int {
	if colour, ok := p.colour[p.location]; ok {
		return colour
	}
	return black
}

func (p *paintBot) getColourAtCoord(coord coord) int {
	if colour, ok := p.colour[coord]; ok {
		return colour
	}
	return black
}

func (p *paintBot) setColour(colour int) {
	p.colour[p.location] = colour
}

func (p *paintBot) moveLeft() {
	p.facingDirection.rotateAnticlockwise()
	p.moveForward()
}

func (p *paintBot) moveRight() {
	p.facingDirection.rotateClockwise()
	p.moveForward()
}
