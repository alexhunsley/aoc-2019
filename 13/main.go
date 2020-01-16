package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

//0 is an empty tile. No game object appears in this tile.
//1 is a wall tile. Walls are indestructible barriers.
//2 is a block tile. Blocks can be broken by the ball.
//3 is a horizontal paddle tile. The paddle is indestructible.
//4 is a ball tile. The ball moves diagonally and bounces off objects.

const (
	tileEmpty int = iota
	tileWall
	tileBlock
	tilePaddle
	tileBall
)

type tile struct {
	x int
	y int
	tileType int
}

var displayChars = ".#M=O"

func main() {
	//log.SetLevel(log.TraceLevel)
	log.SetLevel(log.DebugLevel)
	code := readFileInput()

	solvePart1(code, true)
	//solvePart1(code, false)
}

func outputGameBoard(outputResults []int, displayTiles *[24][45]int, playerScore *int, showBoard bool) (int, int) {
	numTiles := len(outputResults) / 3

	paddleX := -1
	ballX := -1

	for i := 0; i < numTiles; i++ {
		x := (outputResults)[i*3]
		y := (outputResults)[i*3+1]

		if x == -1 {
			// it's a score
			*playerScore = (outputResults)[i*3+2]
			continue
		}
		tileType := (outputResults)[i*3+2]
		tile := tile{x, y, tileType}

		if tileType == tilePaddle {
			paddleX = x
		} else if tileType == tileBall {
			ballX = x
		}
		displayTiles[tile.y][tile.x] = tile.tileType
	}

	///// output the board
	if showBoard {
		for y := 0; y < 24; y++ {
			for x := 0; x < 45; x++ {
				char := displayChars[displayTiles[y][x]]
				fmt.Print(string(char))
			}
			fmt.Println("")
		}
	}
	fmt.Println("SCORE: ", *playerScore)

	return paddleX, ballX
}

func solvePart1(code []int, useAIToPlay bool) {
	// outputs from the intcode machine (tiles, score)
	var outputVals []int
	displayTiles := [24][45]int{}
	playerScore := 0

	machine := IntMachine{
		code:           &code,
		programCounter: 0,
		sparseMemory:   &(map[int]int{}),
	}

	ballX := -1

	inputCallback := func() int {
		//output the map and score, return a paddle dirn
		paddleX, newBallX := outputGameBoard(outputVals, &displayTiles, &playerScore, !useAIToPlay)

		outputVals = []int{}

		if useAIToPlay {
			dir := 0

			ballX = newBallX
			paddleBallDelta := paddleX - ballX

			if paddleBallDelta > 0 {
				dir = -1
			} else if paddleBallDelta < 0 {
				dir = 1
			} else {
				dir = 0
			}
			return dir
		}

		i := getPaddleInputFromUser()
		return i - 2
	}

	outputCallback := func(outputValue int) {
		outputVals = append(outputVals, outputValue)
	}

	// this first run will do breakout machine setup then HALT to wait for quarters to be inserted
	runcode(&machine, inputCallback, outputCallback)

	outputGameBoard(outputVals, &displayTiles, &playerScore, !useAIToPlay)

	// insert 2 quarters after machine setup and showing first board
	machine.poke(0, 2)

	// this call will loop until game over or you win - either way,
	// it will halt when done
	runcode(&machine, inputCallback, outputCallback)

	// need to paint board one last time to see the final score
	outputGameBoard(outputVals, &displayTiles, &playerScore, true)

	fmt.Println("Day 13 part 2 solution: ", playerScore)
}

// get input of 1, 2, or 3 from user (for left, stay, right respectively)
func getPaddleInputFromUser() int {
	userInput := ""
	for len(userInput) == 0 || (len(userInput) == 1 && (userInput[0] < '1' || userInput[0] > '3')) {
		fmt.Print("Move paddle: 1 - left, 2 - stay, 3 - right: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userInput = scanner.Text()
	}
	i, _ := strconv.Atoi(userInput)
	return i
}

func readFileInput() []int {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

	codeAsStrings := strings.Split(line, ",")

	var code = []int{}

	for _, i := range codeAsStrings {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		code = append(code, j)
	}
	return code
}
