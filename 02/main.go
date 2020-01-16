package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Opcode uint

const (
	ADD Opcode = 1
	MULT = 2
	END = 99
)

func main() {
	fmt.Println("ji ", MULT)

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	// Default scanner is bufio.ScanLines. Lets use ScanWords.
	scanner.Split(bufio.ScanWords)

	success := scanner.Scan()
	if success == false {
		// False on error or EOF. Check error
		err = scanner.Err()
		if err == nil {
			// reached EOF
			//log.Println("Scan completed and reached EOF")
			log.Fatal("Didn't find a line in input.txt")
		} else {
			log.Fatal(err)
		}
	}
	var line = scanner.Text()
	fmt.Println("line = ", line)

	// Split on comma.
	codeAsStrings := strings.Split(line, ",")

	var codeInput = []int{}

	for _, i := range codeAsStrings {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		codeInput = append(codeInput, j)
	}

	fmt.Println(codeInput)

	outerloop:
	for verb := 0; verb < 100; verb++ {
		for noun := 0; noun < 100; noun++ {
			//fmt.Println("==================== verb. noun =", verb, noun)
			//fmt.Println(codeInput)

			// doesn't work! codeInput is a slice!
			// code := codeInput

			code := make([]int, len(codeInput))
			copy(code, codeInput)

			codeAtZeroAddress := runcode(code, verb, noun)

			if noun == 12 && verb == 2 {
				fmt.Println("Part 1 answer: ", codeAtZeroAddress)
			}

			if codeAtZeroAddress == 19690720 {
				part2Answer := 100 * noun + verb
				fmt.Println("Part 2 answer: ", part2Answer)
				break outerloop
			}
		}
	}
}

func runcode(code []int, verb int, noun int) int {
	code[1] = noun
	code[2] = verb

	codeIndex := 0
	// Display all elements.
	for true {
		//fmt.Println(code[codeIndex])

		opcode := Opcode(code[codeIndex])

		switch opcode {
		case ADD:
			addr1 := code[codeIndex + 1]
			addr2 := code[codeIndex + 2]
			dest := code[codeIndex + 3]
			//fmt.Printf("+ %d, %d --> %d\n", addr1, addr2, dest)

			code[dest] = code[addr1] + code[addr2]
			//fmt.Println("    code is now: ", code)
		case MULT:
			addr1 := code[codeIndex + 1]
			addr2 := code[codeIndex + 2]
			dest := code[codeIndex + 3]
			//fmt.Printf("* %d, %d --> %d\n", addr1, addr2, dest)

			code[dest] = code[addr1] * code[addr2]
			//fmt.Println("    code is now: ", code)
		case END:
			//fmt.Printf("HALT\n")
			break
		default:
			log.Fatal("Found unrecognized opcode: ", opcode)
			panic("Abort")
		}
		codeIndex += 4

		if opcode == END {
			//fmt.Println("Halted at code: ", code)
			return code[0]
		}
	}
	return 0
}
