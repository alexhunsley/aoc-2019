package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type SpaceObject struct {
	parentObject *SpaceObject
	name string
	orbitDepth int
}

func (s SpaceObject) String() string {
	var parentName string = "<none>"

	if s.parentObject != nil {
		parentName = s.parentObject.name
	}
	return fmt.Sprintf("[%s, orbDepth = %d, parName = %s]", s.name, s.orbitDepth, parentName)
}

var spaceObjects = map[string]SpaceObject{}

func findSpaceObjectForName(objName string) SpaceObject {
	val, ok := spaceObjects[objName]
	if !ok {
		panic("Didn't find a space object when I expected one")
	}
	return val
}

func findOrCreateSpaceObjectForName(objName string, parentSpaceObject *SpaceObject, orbitDepth int) SpaceObject {
	val, ok := spaceObjects[objName]
	if !ok {
		newObj := SpaceObject{name: objName}
		newObj.orbitDepth = orbitDepth
		if parentSpaceObject != nil {
			newObj.parentObject = parentSpaceObject
		}

		spaceObjects[objName] = newObj
		return newObj
	}
	return val
}

// get the values array for a key in orbit map, creating the array
// if it doesn't already exit
func valuesArrayForOrbitMapEntry(orbitsMap map[string][]string, name string) []string {
	namesArray, ok := orbitsMap[name]
	if !ok {
		namesArray = []string{}
		orbitsMap[name] = namesArray
	}
	return namesArray
}

func readOrbitsFile() map[string][]string {
	// map from string to (string array)
	// e.g. ABC -> [DEF, GHI],
	//      XYZ -> [PQR]
	var orbitsMap = map[string][]string{}

	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ")")

		namesArray := valuesArrayForOrbitMapEntry(orbitsMap, parts[0])
		orbitsMap[parts[0]] = append(namesArray, parts[1])
	}
	file.Close()

	return orbitsMap
}

func calcOrbits(orbitsMap map[string][]string, bodyName string, parentSpaceObject *SpaceObject, orbitDepth int) *SpaceObject {
	spaceObject := findOrCreateSpaceObjectForName(bodyName, parentSpaceObject, orbitDepth)

	childBodies := orbitsMap[bodyName]

	for _, childBodyName := range childBodies {
		calcOrbits(orbitsMap, childBodyName, &spaceObject, orbitDepth + 1)
	}
	return &spaceObject
}

func sumOrbits() int {
	orbitsTotal := 0

	for _, spaceObj := range spaceObjects {
		orbitsTotal += spaceObj.orbitDepth
	}
	return orbitsTotal
}

func findAncestorBodies(spaceObj *SpaceObject) []*SpaceObject {
	ancestors := []*SpaceObject{}

	currBody := spaceObj

	for currBody != nil {
		ancestors = append(ancestors, currBody)
		currBody = currBody.parentObject
	}
	return ancestors
}

func findFirstCommonAncestor(spaceObjects1 []*SpaceObject, spaceObjects2 []*SpaceObject) *SpaceObject {
	for _, spaceObjIn1 := range spaceObjects1 {
		for _, spaceObjIn2 := range spaceObjects2 {
			if spaceObjIn1.name == spaceObjIn2.name {
				return spaceObjIn1
			}
		}
	}
	panic("No common ancestor - should never happen!")
}

func main() {
	// get map from string to string array of orbit relations.
	// e.g. COM -> [DEF, GHI],
	//      DEF -> [PQR]
	orbitsMap := readOrbitsFile()

	calcOrbits(orbitsMap, "COM", nil, 0)

	totalOrbits := sumOrbits()

	fmt.Println("Part 1 solution: ", totalOrbits)

	// Part2:
	//
	// Strategy: find your ancestor bodies list, and santas.
	// First common item in them identifies the place you have to go via.
	// Then the total item travel dist is (orbitDeptYou - commonOrbitDepth) + (orbitDeptSanta - commonOrbitDepth) - 2
	//    = orbitDeptYou + orbitDeptSanta - 2 * (commonOrbitDepth + 1)
	//
	// NOTE: we have to subtract 2 as per instructions - see bit about we end up orbitting same thing as santa
	//

	youBody := findSpaceObjectForName("YOU")
	santaBody := findSpaceObjectForName("SAN")

	yourAncestors := findAncestorBodies(youBody.parentObject)
	santaAncestors := findAncestorBodies(santaBody.parentObject)

	firstCommonAncestor := findFirstCommonAncestor(yourAncestors, santaAncestors)

	travelDistance := youBody.orbitDepth + santaBody.orbitDepth - 2 * (firstCommonAncestor.orbitDepth + 1)

	fmt.Println("Part 2 solution: ", travelDistance)

}
