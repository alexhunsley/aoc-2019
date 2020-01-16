package main

import (
	"bufio"
	"math"
	"os"
)

// In effect, in Day 11 we are implementing a 2D tape in a turing machine.
//

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	black int = iota
	white
)

//////////////////////////////////////////////////////////////////////////////////////////////////////
// main

func main() {
	os.Remove("go.log")

	log.SetLevel(log.InfoLevel)

	logFile, logErr := os.OpenFile("go.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if logErr == nil {
		log.SetOutput(logFile)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	//////////////// Part 1

	bot := paintBot{colour: make(map[coord]int), facingDirection: compassDirection{xPart: 0, yPart: -1}}

	line := readInputFile()

	numTiles := solvePart1(line, &bot)
	fmt.Println("Part 1 solution: ", numTiles)

	//////////////// Part 2

	bot = paintBot{colour: make(map[coord]int), facingDirection: compassDirection{xPart: 0, yPart: -1}}
	bot.setColour(white)

	solvePart1(line, &bot)

	// gives AHLCPRAL
	solvePart2(&bot)
}

func solvePart2(bot *paintBot) {
	topLeft := coord{
		x: math.MaxInt64,
		y: math.MaxInt64,
	}
	bottomRight := coord{
		x: math.MinInt64,
		y: math.MinInt64,
	}

	// find range of the painted tiles
	for tileCoord, _ := range bot.colour {
		if tileCoord.x < topLeft.x {
			topLeft.x = tileCoord.x
		}
		if tileCoord.y < topLeft.y {
			topLeft.y = tileCoord.y
		}
		if tileCoord.x > bottomRight.x {
			bottomRight.x = tileCoord.x
		}
		if tileCoord.y > bottomRight.y {
			bottomRight.y = tileCoord.y
		}
	}

	for y := topLeft.y; y <= bottomRight.y; y += 1 {
		for x := topLeft.x; x <= bottomRight.x; x += 1 {
			col := bot.getColourAtCoord(coord{
				x: x,
				y: y,
			})

			fmt.Print([]string{" ", "*"}[col])
		}
		fmt.Println("")
	}
}

func solvePart1(line string, bot *paintBot) int {
	codeAsStrings := strings.Split(line, ",")

	var code = []int{}

	for _, i := range codeAsStrings {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		code = append(code, j)
	}

	part1Code := make([]int, len(code))
	copy(part1Code, code)

	inputChan := make(chan int)
	outputChan := make(chan int)

	machine := IntMachine{
		code:           &part1Code,
		programCounter: 0,
		sparseMemory:   &(map[int]int{}),
	}

	go runcode(&machine, inputChan, outputChan)

	// This loop is a little brittle in that it assumes the intcode computer program will strictly follow
	// a 'get input, do 2 outputs, repeat' pattern.
	// Happily, things are indeed this way, so it works.
	// But if the intcode program asked for the input several times between sending output values,
	// we'd get deadlock here and have to use perhaps a callback for input instead of a blocking inputChannel fetch.
	for true {
		// to avoid deadlock, whereby intmachine is waiting on sending -1 via outputChan (to signal halt),
		// while we are waiting on putting something into the inputChan here.
		go func() {
			inputChan <- bot.getColour()
		}()

		colourToPaint := <-outputChan

		if colourToPaint < 0 {
			// signals halt
			return len(bot.colour)
		}

		bot.setColour(colourToPaint)
		turnDirection := <-outputChan

		if turnDirection == 0 {
			bot.moveLeft()
		} else {
			bot.moveRight()
		}
	}
	// shouldn't get here
	return 0
}

func readInputFile() string {
	file, err := os.Open("input.txt")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)

	success := scanner.Scan()
	if success == false {
		// False on error or EOF. Check error
		err = scanner.Err()
		if err == nil {
			log.Fatal("Didn't find a line in input.txt")
		} else {
			log.Fatal(err)
		}
	}
	var line = scanner.Text()
	return line
}
