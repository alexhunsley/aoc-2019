package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// we're reading then processing line-by-line. Obviously this might not be good
// for huge inputs.
func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	// Default scanner is bufio.ScanLines. Lets use ScanWords.
	scanner.Split(bufio.ScanWords)

	var totalFuelNeededPart1 int64 = 0
	var totalFuelNeededPart2 int64 = 0

	// Scan for next token.
	for true {
		success := scanner.Scan()
		if success == false {
			// False on error or EOF. Check error
			err = scanner.Err()
			if err == nil {
				// reached EOF
				//log.Println("Scan completed and reached EOF")
				break
			} else {
				log.Fatal(err)
			}
		}

		var massStr = scanner.Text()

		mass, err := strconv.ParseInt(massStr, 10, 64)
		if err != nil {
			panic(err)
		}
		//fmt.Println("Read a mass: ", mass)

		var fuelNeededForMainMass = calcfuel(mass)

		totalFuelNeededPart1 += fuelNeededForMainMass

		totalFuelNeededPart2 += calcfuelrecursive(mass)
	}
	fmt.Println("Part 1 Total fuel needed: ", totalFuelNeededPart1)
	fmt.Println("Part 2 Total fuel needed: ", totalFuelNeededPart2)
}

/*Fuel required to launch a given module is based on its mass. Specifically, to find the fuel required for a module,
take its mass, divide by three, round down, and subtract 2.
*/
func calcfuel(mass int64) int64 {
	// integer division in golang rounds down for us
	return mass / 3  - 2
}

func calcfuelrecursive(mass int64) int64 {
	var fuelTotal int64 = 0

	for true {
		var fuelNeeded = calcfuel(mass)
		if fuelNeeded <= 0 {
			//fmt.Println("recurse, HIT 0 OR LESS! exiting.  fuel = ", fuelNeeded)
			break
		}
		//fmt.Println("recurse, fuel needed: ", fuelNeeded)
		fuelTotal += fuelNeeded
		mass = fuelNeeded
	}
	return fuelTotal
}