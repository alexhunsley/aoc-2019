package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Assumptions about this problem:
//  1. wires only intersect at right angles - 2 wires cannot go in same direction at overlapping coordinates.
//  2. a wire doesn't intersect itself (so we only need to check for intersections with the other wire)
//
// These assumptions mean we only have to check for horizontal wire sections from wire A hitting vertical
// sections on wire B, and vice versa.

import (
	"fmt"
	"strconv"
)

type LineType uint

const (
	HORIZ LineType = iota
	VERT
)

type coord struct {
	x int
	y int
}

func (c coord) manhattanModulus() int {
	return abs(c.x) + abs(c.y)
}

type lineDef struct {
	direction LineType
	// the Y coord of a horizontal line, or the X coordinate of a vertical line
	staticCoordinate int
	// start and end X coordinates for horiz line, or similarly Y coords for a vertical line.
	// the start value is always less than the end value. this is to make intersection logic simpler.
	otherCoordinateStart int
	otherCoordinateEnd int
	// true if the 'flow' direction of the wire is opposite to the start, end order
	startEndCoordsFlipped bool
	// cumulative grid distance from start of the wire to this line section
	distToStart int
}

func (l lineDef) pr() string {
	var dirnStr string

	if l.direction == HORIZ {
		dirnStr = "HORIZ"
	} else {
		dirnStr = "VERT"
	}
	return fmt.Sprintf("(%s, static coord: %d, other start: %d, other end: %d, flipped: %s, cum dist: %d)", dirnStr, l.staticCoordinate, l.otherCoordinateStart, l.otherCoordinateEnd, l.startEndCoordsFlipped, l.distToStart)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)

	wireLine1, err := reader.ReadString('\n')
	wireLine1 = strings.TrimSuffix(wireLine1, "\n")

	if err != nil {
		panic("Couldn't read wire line 1 from input")
	}

	wireLine2, err := reader.ReadString('\n')
	wireLine2 = strings.TrimSuffix(wireLine2, "\n")

	if err != nil {
		panic("Couldn't read wire line 2 from input")
	}

	wireLineParts1 := parseWireLine(wireLine1)
	wireLineParts2 := parseWireLine(wireLine2)

	_, closestIntersection, smallestTotalDist := findIntersections(wireLineParts1, wireLineParts2)

	fmt.Println("Part 1: closest intersection, dist (ANSWER!) = ", closestIntersection, closestIntersection.manhattanModulus())
	fmt.Println("Part 2: smallest total dist = ", smallestTotalDist)
}

func findIntersections(wireLineParts1 []lineDef, wireLineParts2 []lineDef) ([]coord, coord, int) {
	horizLineTest := func(linedef lineDef) bool { return linedef.direction == HORIZ }
	vertLineTest := func(linedef lineDef) bool { return linedef.direction == VERT }

	horizLines1 := filterLineDefs(wireLineParts1, horizLineTest)
	vertLines1 := filterLineDefs(wireLineParts1, vertLineTest)

	horizLines2 := filterLineDefs(wireLineParts2, horizLineTest)
	vertLines2 := filterLineDefs(wireLineParts2, vertLineTest)

	var intersectionCoords = []coord{}

	var closestIntersectionCoord coord
	var closestIntersectionDist int = -1

	smallestTotalDist := -1

	// intersection of horiz lines in wire 1 with all vert lines in wire 2.
	// This and the following block can easily be factored into a single function that we call twice,
	// but I can't be bothered today!
	for _, hLineInWire1 := range horizLines1 {
		for _, vLineInWire2 := range vertLines2 {
			if hLineInWire1.otherCoordinateStart < vLineInWire2.staticCoordinate &&
				hLineInWire1.otherCoordinateEnd > vLineInWire2.staticCoordinate &&
				vLineInWire2.otherCoordinateStart < hLineInWire1.staticCoordinate &&
				vLineInWire2.otherCoordinateEnd > hLineInWire1.staticCoordinate {

				intersectionCoord := coord{x: vLineInWire2.staticCoordinate, y: hLineInWire1.staticCoordinate}

				// calc distance along both wires that intersect
				var distHLine int
				if hLineInWire1.startEndCoordsFlipped {
					distHLine = hLineInWire1.distToStart + abs(hLineInWire1.otherCoordinateEnd - vLineInWire2.staticCoordinate)
				} else {
					distHLine = hLineInWire1.distToStart + abs(hLineInWire1.otherCoordinateStart - vLineInWire2.staticCoordinate)
				}

				dist := intersectionCoord.manhattanModulus()
				if closestIntersectionDist < 0 || dist < closestIntersectionDist {
					closestIntersectionDist = dist
					closestIntersectionCoord = intersectionCoord
				}
				intersectionCoords = append(intersectionCoords, intersectionCoord)

				var distVLine int
				if vLineInWire2.startEndCoordsFlipped {
					distVLine = vLineInWire2.distToStart + abs(vLineInWire2.otherCoordinateEnd - hLineInWire1.staticCoordinate)
				} else {
					distVLine = vLineInWire2.distToStart + abs(vLineInWire2.otherCoordinateStart - hLineInWire1.staticCoordinate)
				}

				totalDist := distHLine + distVLine

				if smallestTotalDist < 0 || totalDist < smallestTotalDist {
					smallestTotalDist = totalDist
				}
			}
		}
	}

	// intersection of vert lines in wire 1 with all horiz lines in wire 2
	for _, vLineInWire1 := range vertLines1 {
		for _, hLineInWire2 := range horizLines2 {
			if vLineInWire1.otherCoordinateStart < hLineInWire2.staticCoordinate &&
				vLineInWire1.otherCoordinateEnd > hLineInWire2.staticCoordinate &&
				hLineInWire2.otherCoordinateStart < vLineInWire1.staticCoordinate &&
				hLineInWire2.otherCoordinateEnd > vLineInWire1.staticCoordinate {

				intersectionCoord := coord{x: vLineInWire1.staticCoordinate, y: hLineInWire2.staticCoordinate}

				// calc distance along both wires that intersect
				var distHLine int
				if vLineInWire1.startEndCoordsFlipped {
					distHLine = vLineInWire1.distToStart + abs(vLineInWire1.otherCoordinateEnd - hLineInWire2.staticCoordinate)
				} else {
					distHLine = vLineInWire1.distToStart + abs(vLineInWire1.otherCoordinateStart - hLineInWire2.staticCoordinate)
				}

				dist := intersectionCoord.manhattanModulus()
				if closestIntersectionDist < 0 || dist < closestIntersectionDist {
					closestIntersectionDist = dist
					closestIntersectionCoord = intersectionCoord
				}
				intersectionCoords = append(intersectionCoords, intersectionCoord)

				var distVLine int
				if hLineInWire2.startEndCoordsFlipped {
					distVLine = hLineInWire2.distToStart + abs(hLineInWire2.otherCoordinateEnd - vLineInWire1.staticCoordinate)
				} else {
					distVLine = hLineInWire2.distToStart + abs(hLineInWire2.otherCoordinateStart - vLineInWire1.staticCoordinate)
				}

				totalDist := distHLine + distVLine

				if smallestTotalDist < 0 || totalDist < smallestTotalDist {
					smallestTotalDist = totalDist
				}
			}
		}
	}
	return intersectionCoords, closestIntersectionCoord, smallestTotalDist
}

func parseWireLine(line string) []lineDef {
	wireMoves := strings.Split(line, ",")

	var wireX int = 0
	var wireY int = 0

	var lines []lineDef

	cumulativeDistToStart := 0

	for _, i := range wireMoves {
		direction := i[0]
		coordAsInt, _ := strconv.Atoi(i[1:])

		coord := coordAsInt

		var newLine lineDef

		switch direction {
		case 'L':
			newLine = lineDef{direction: HORIZ, staticCoordinate: wireY, otherCoordinateStart: wireX - coord, otherCoordinateEnd: wireX, startEndCoordsFlipped: true, distToStart: cumulativeDistToStart }
			cumulativeDistToStart += coord
			wireX = wireX - coord
		case 'R':
			//fmt.Println("got L or R")
			newLine = lineDef{direction: HORIZ, staticCoordinate: wireY, otherCoordinateStart: wireX, otherCoordinateEnd: wireX + coord, startEndCoordsFlipped: false, distToStart: cumulativeDistToStart }
			cumulativeDistToStart += coord
			wireX = wireX + coord
		case 'U':
			newLine = lineDef{direction: VERT, staticCoordinate: wireX, otherCoordinateStart: wireY, otherCoordinateEnd: wireY + coord, startEndCoordsFlipped: false, distToStart: cumulativeDistToStart }
			cumulativeDistToStart += coord
			wireY = wireY + coord
		case 'D':
			//fmt.Println("got U or D")
			newLine = lineDef{direction: VERT, staticCoordinate: wireX, otherCoordinateStart: wireY - coord, otherCoordinateEnd: wireY, startEndCoordsFlipped: true, distToStart: cumulativeDistToStart }
			cumulativeDistToStart += coord
			wireY = wireY - coord
		default:
			panic("Got unexpected char for direction that wasn't L, R, U or D")
		}

		lines = append(lines, newLine)
	}
	return lines
}

func filterLineDefs(lineDefs []lineDef, test func(lineDef) bool) (ret []lineDef) {
	for _, s := range lineDefs {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
