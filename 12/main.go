package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// for set impl
type void struct{}
var member void

type vec3 struct {
	x int
	y int
	z int
}

type moon struct {
	pos vec3
	vel vec3
}

func (v1 vec3) add(v2 vec3) vec3 {
	return vec3{v1.x + v2.x, v1.y + v2.y, v1.z + v2.z}
}

func (v1 vec3) sub(v2 vec3) vec3 {
	return vec3{v1.x - v2.x, v1.y - v2.y, v1.z - v2.z}
}

// return an int array containing the pos then vec components
func (m moon) collapsedData() []int {
	//return []int{m.pos.x, m.pos.y, m.pos.z, m.vel.x, m.vel.y, m.vel.z}
	return []int{m.pos.x, m.pos.y, m.pos.z, m.vel.x, m.vel.y, m.vel.z}
}

func (v1 vec3) vecSgn() vec3 {
	return vec3{sgnWithZero(v1.x), sgnWithZero(v1.y), sgnWithZero(v1.z)}
}

func (m1 moon) updatedPosition() vec3 {
	return m1.pos.add(m1.vel)
}

//func (m1 *moon) updatePosition() {
//	m1.pos = m1.pos.add(m1.vel)
//}

func main() {
	lines := readInput()
	moonData := parseInitialPositions(lines)

	fmt.Println("moonData after reading: ", moonData)
	simulateTimeSteps(1000000000000, moonData)

	totalEnergy := calcTotalEnergy(moonData)

	fmt.Println("Day 12 part 1 solution: ", totalEnergy)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func calcTotalEnergy(moonData []moon) interface{} {
	numMoons := len(moonData)

	totalEnergy := 0

	for m := 0; m < numMoons; m++ {
		moon := moonData[m]
		potentialEnergy := abs(moon.pos.x) + abs(moon.pos.y) + abs(moon.pos.z)
		kineticEnergy := abs(moon.vel.x) + abs(moon.vel.y) + abs(moon.vel.z)

		totalEnergy += potentialEnergy * kineticEnergy
	}
	return totalEnergy
}

func simulateTimeSteps(numSteps int, moonData []moon) {
	numMoons := len(moonData)

	fmt.Println("got initial data: ", moonData)

	var initialData = make([]moon, len(moonData))
	copy(initialData, moonData)

	//var initialVels = [][3]int{moonData[0].vel, moonData[1].vel, moonData[2].vel, moonData[3].vel}

	fmt.Println("initial data: ", initialData)

	moonRepeatPeriod := [4]int{-1, -1, -1, -1}
	moonRepeatsFound := 0

	// empty set for seen coords
	//previousMoonData := make(map[string]int)
	//previousMoonData := [4]map[string]int{}
	previousMoonData := make([]map[string]int, 4)
	for ii := 0; ii < 4; ii++ {
		previousMoonData[ii] = make(map[string]int)
	}

	//set["Foo"] = member

	for i := 1; i < numSteps+1; i++ {
		//fmt.Println("============================== step ", i, " moon[0] = ", moonData[0])

		// update gravity
		for m1 := 0; m1 < numMoons-1; m1++ {
			for m2 := m1 + 1; m2 < numMoons; m2++ {
				deltaSignVel := moonData[m2].pos.sub(moonData[m1].pos).vecSgn()

				//fmt.Println("delta sign vel: ", deltaSignVel)
				moonData[m1].vel = moonData[m1].vel.add(deltaSignVel)
				moonData[m2].vel = moonData[m2].vel.sub(deltaSignVel)
			}
		}

		//fmt.Println("velocities now: ")
		//for _, moon := range moonData {
		//	fmt.Println(moon.vel)
		//}

		// update positions
		for j, m := range moonData {
			moonData[j].pos = m.updatedPosition()
		}

		if i % 1000000 < 1 {
			fmt.Println("positions, vel now at step: ", i)
			for _, moon := range moonData {
				fmt.Println(moon.pos, moon.vel)
			}
			fmt.Println("num hashes: ", len(previousMoonData)) //, " set = ", previousMoonData)
		}

		//for moonIdxForHistoryCheck := 0; moonIdxForHistoryCheck < numMoons; moonIdxForHistoryCheck++ {
		for moonIdxForHistoryCheck := 0; moonIdxForHistoryCheck < numMoons; moonIdxForHistoryCheck++ {
			if moonRepeatPeriod[moonIdxForHistoryCheck] > 0 {
				continue
			}

			dataKey := computeStringHashForList(moonData[moonIdxForHistoryCheck].collapsedData())
			oldIndex, exists := previousMoonData[moonIdxForHistoryCheck][dataKey]
 
			if exists {
				moonRepeatsFound += 1

				moonRepeatPeriod[moonIdxForHistoryCheck] = i - oldIndex

				fmt.Println("Found old state repeated! moonIdx, i = ", moonIdxForHistoryCheck, i, " oldIndex is ", oldIndex, "   period is ", moonRepeatPeriod[moonIdxForHistoryCheck])
				fmt.Println("Repeated data = ", moonData[moonIdxForHistoryCheck])
				if moonRepeatsFound == 4 {
					return
				}
			}

			previousMoonData[moonIdxForHistoryCheck][dataKey] = i
		}
		//fmt.Println("positions now at step: ", i)
		//for _, moon := range moonData {
		//	fmt.Println(moon)
		//}

		//if i > 0 && (testEq(initialCoords[0], moonData[0].pos) && testEq(initialCoords[1], moonData[1].pos) &&
		//		testEq(initialCoords[2], moonData[2].pos) && testEq(initialCoords[3], moonData[3].pos) &&
		//	testEq(initialVels[0], moonData[0].vel) && testEq(initialVels[1], moonData[1].vel) &&
		//	testEq(initialVels[2], moonData[2].vel) && testEq(initialVels[2], moonData[2].vel)) {

		//if i > 0 && (testEq(initialCoords[0], moonData[0].pos) &&
		//	testEq(initialVels[0], moonData[0].vel)) {
		//if i > 0 && (testEq(initialCoords[1], moonData[1].pos) &&
		//	testEq(initialVels[1], moonData[1].vel)) {

		//coordIdx := 1
		//
		//if i > 0 && (testEq(initialCoords[coordIdx], moonData[coordIdx].pos) &&
		//	testEq(initialVels[coordIdx], moonData[coordIdx].vel)) {
		//// just look for 0 velocity - half way point
		////if i > 0 && testEq(initialVels[coordIdx], moonData[coordIdx].vel) {
		//
		//	fmt.Println("Part 2 complete! num steps = ", i, " 2 x steps = ", 2*i, "  initCoord[0], moonDataPos[0] = ", initialCoords, moonData[0].pos, moonData[1].pos, moonData[2].pos, moonData[3].pos)
		//	return
		//}
	}

	// bugger. it's the "matching a point previously in time" thing, isn't it?
	// I need to match anything in the past, not just the start state!

	fmt.Println("final coords = ", moonData)
	fmt.Println("initialCoords = ", initialData)
}

func testEq(a, b [3]int) bool {

	// If one is nil, the other must also be nil.
	//if (a == nil) != (b == nil) {
	//	return false;
	//}

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func sgnWithZero(x int) int {
	if x == 0 {
		return 0
	}
	if x < 0 {
		return -1
	}
	return 1
}

func parseInitialPositions(lines []string) []moon {
	var moons []moon

	re1 := regexp.MustCompile("<x=(-?\\d*), y=(-?\\d*), z=(-?\\d*)>")

	for _, line := range lines {
		matches := re1.FindStringSubmatch(line)

		xx, _ := strconv.Atoi(matches[1])
		yy, _ := strconv.Atoi(matches[2])
		zz, _ := strconv.Atoi(matches[3])

		moon := moon{
			pos: vec3{xx, yy, zz},
			vel: vec3{0, 0, 0},
		}
		moons = append(moons, moon)
	}
	return moons
}

func readInput() []string {
	file, err := os.Open("input3.txt")

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
