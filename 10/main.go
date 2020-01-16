package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type VisibilityResult struct {
	satelliteX int
	satelliteY int
	numSatsVisible int
}

type polarCoordinate struct {
	angle float64
	dist float64
	x int
	y int
}

func main() {
	lines := readInput()

	solvePart1(lines)

	stationX := 22
	stationY := 25

	solvePart2(lines, stationX, stationY)
}

func readInput() []string {
	file, err := os.Open("input.txt")

	if err != nil {
		panic("failed opening input file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines = []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func solvePart2(lines []string, stationX int, stationY int) {
	width := len(lines[0])
	height := len(lines)

	polarCoords := calcPolarCoordinates(height, width, lines, stationX, stationY)

	sortPolarCoords(polarCoords)

	destroyAsteroidsInOrder(polarCoords)
}

func destroyAsteroidsInOrder(polarCoords []polarCoordinate) {
	asteroidDestructionIndex := 0

	var lastMatchedAngle float64 = -1000

	asteroidIndex := 0
	var destroyedAsteroid polarCoordinate

	for asteroidDestructionIndex < 200 {
		for true {
			angleDiff := math.Abs(polarCoords[asteroidIndex].angle - lastMatchedAngle)

			if angleDiff > 0.00000001 {
				//// N.B. we can get stuck on the last asteroid(s) if they are all at the same angle!
				break
			}
			asteroidIndex = (asteroidIndex + 1) % len(polarCoords)
		}
		asteroidDestructionIndex += 1

		lastMatchedAngle = polarCoords[asteroidIndex].angle

		destroyedAsteroid = polarCoords[asteroidIndex]

		// remove the destroyed asteroid
		resultPolarCoords := polarCoords[:asteroidIndex]
		polarCoords = append(resultPolarCoords, polarCoords[asteroidIndex+1:]...)

		// check if our asteroidIndex is now out of bounds!
		if asteroidIndex == len(polarCoords) {
			asteroidIndex = 0
		}
	}
	fmt.Println("Part 2 solution: asteroid = ", destroyedAsteroid)
}

func sortPolarCoords(polarCoords []polarCoordinate) {
	// sort the polar coords majorly by increasing angle and minorly by increasing distance (radius)
	sort.SliceStable(polarCoords, func(i int, j int) bool {
		p1 := polarCoords[i]
		p2 := polarCoords[j]

		if math.Abs(p1.angle-p2.angle) < 0.00000001 {
			return p1.dist < p2.dist
		}
		if p1.angle < p2.angle {
			return true
		}
		return false
	})
}

func calcPolarCoordinates(height int, width int, lines []string, stationX int, stationY int) []polarCoordinate {
	polarCoords := []polarCoordinate{}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if lines[y][x] == '.' {
				continue
			}
			if x == stationX && y == stationY {
				continue
			}

			dx := x - stationX
			dy := y - stationY

			angle := math.Atan2(float64(dx), float64(-dy))
			if angle < 0 {
				angle += 2 * math.Pi
			}
			// don't need to calc the sqrt since we're just comparing distances
			polarCoord := polarCoordinate{dist: float64(dx*dx + dy*dy), angle: angle, x: x, y: y}

			polarCoords = append(polarCoords, polarCoord)
		}
	}
	return polarCoords
}

//========================================================================================================

func solvePart1(lines []string) {
	width := len(lines[0])
	height := len(lines)

	visResult := VisibilityResult{}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if lines[y][x] == '.' {
				continue
			}

			numSatsVisible := calcSatsVisible(lines, x, y)
			if visResult.numSatsVisible == -1 || numSatsVisible > visResult.numSatsVisible {
				visResult.satelliteX = x
				visResult.satelliteY = y
				visResult.numSatsVisible = numSatsVisible
			}
		}
	}
	fmt.Println("Part 1 solution: max of ", visResult.numSatsVisible, " (at xy = ", visResult.satelliteX, visResult.satelliteY, ")")
}

// return the number of satellites visible from the given coord.
func calcSatsVisible(lines []string, satelliteX int, satelliteY int) int {
	width := len(lines[0])
	height := len(lines)

	satellitesVisible := 0

	for y := 0; y < height; y++ {
		out:
		for x := 0; x < width; x++ {
			if lines[y][x] != '#' || (x == satelliteX && y == satelliteY) {
				continue
			}

			dx := x - satelliteX
			dy := y - satelliteY

			var jumpX int
			var jumpY int

			// find our minimum 'jump distance' to scan grid in the hunt for visibility blockers
			if dx == 0 {
				jumpX = 0
				jumpY = sgn(dy)
			} else if dy == 0 {
				jumpX = sgn(dx)
				jumpY = 0
			} else {
				gcd := GCD(abs(dy), abs(dx))

				if gcd == 1 {
					// we have an irreducible fraction for the delta,
					// so there can't be any other satellites in the way
					satellitesVisible += 1
					continue
				}
				jumpX = dx / gcd
				jumpY = dy / gcd
			}

			examineX := satelliteX + jumpX
			examineY := satelliteY + jumpY

			for examineX != x || examineY != y {

				if lines[examineY][examineX] == '#' {
					continue out
				}
				examineX += jumpX
				examineY += jumpY
			}
			// no blocking satellites found!
			satellitesVisible += 1
		}
	}
	return satellitesVisible
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sgn(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}
