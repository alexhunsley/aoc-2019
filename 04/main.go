package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	start := time.Now()

	//startCode := 156218
	//endCode := 652527

	// Use easily calculable clamps to start and end values to reduce runtime.
	// (easy to write a func to do this, but I did these by hand).
	// When I do this on iMac Pro, the time per run is 0.02685
	// instead of 0.03037, which is a saving of 12% on naive start, end values.
	startCode := 156666
	endCode := 599999

	numValidCodesPart1 := 0
	numValidCodesPart2 := 0

	numRuns := 100

	for run := 0; run < numRuns; run++ {
		for code := startCode; code <= endCode; code++ {
			if isValidNumberPart1(strconv.Itoa(code)) {
				numValidCodesPart1 += 1
			}
			if isValidNumberPart2(strconv.Itoa(code)) {
				numValidCodesPart2 += 1
			}
		}
	}
	elapsed := time.Since(start)
	log.Printf("Solution took %s", float64(elapsed) / float64(numRuns) / float64(time.Second))

	fmt.Println("Part 1: Num valid codes: ", numValidCodesPart1)
	fmt.Println("Part 2: Num valid codes: ", numValidCodesPart2)
}

//Criteria for valid number:
//6 digits long.
//Two adjacent digits are the same (like 22 in 122345).
//Going from left to right, the digits never decrease; they only ever increase or stay the same (like 111123 or 135679).
func isValidNumberPart1(number string) bool {
	if len(number) != 6 {
		return false
	}

	haveSeenRepeatedDigit := false

	for pos := 0; pos < 5; pos++ {
		if number[pos] == number[pos + 1] {
			haveSeenRepeatedDigit = true
		} else if number[pos] > number[pos + 1] {
			return false
		}
	}
	return haveSeenRepeatedDigit
}

//Criteria for valid number:
//6 digits long.
///Two and only two adjacent digits are the same (like 22 in 223333), larger groupings don't count.
//Going from left to right, the digits never decrease; they only ever increase or stay the same (like 111123 or 135679).
func isValidNumberPart2(number string) bool {
	if len(number) != 6 {
		return false
	}

	digitForCurrentGroup := uint8(0)
	groupingCount := 0

	haveSeenSingleRepeatedDigit := false

	for pos := 0; pos < 5; pos++ {
		if number[pos] == number[pos + 1] {
			if digitForCurrentGroup == uint8(0) || digitForCurrentGroup == number[pos] {
				groupingCount += 1
			}
		} else {
			if number[pos] > number[pos + 1] {
				return false
			}
			// We've broken a repeated digit sequence.
			// Check for a digit that occurred twice un a row.
			if groupingCount == 1 {
				// can't return yet, since we need to catch invalid codes that decrease a value
				haveSeenSingleRepeatedDigit = true
			}
			groupingCount = 0
			digitForCurrentGroup = uint8(0)
		}
	}

	// check for final pair of digits being the same
	if groupingCount == 1 {
		return true
	}
	return haveSeenSingleRepeatedDigit
}
