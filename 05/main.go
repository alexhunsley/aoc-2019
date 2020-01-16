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
	ADD Opcode = iota + 1
	MULT
	INP
	OUT
	JIT
	JIF
	LT
	EQ
	END = 99
)

var instrLengths = []int{-1,
	4, //ADD
	4, //MULT
	2, //INP
	2, //OUT
	3, //JIT
	3, //JIF
	4, //LT
	4, //EQ
	1}

func padInstruction(instrCode int) string {
	return fmt.Sprintf("%05d", instrCode)
}

func main() {
	file, err := os.Open("input.txt")
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

	codeAsStrings := strings.Split(line, ",")

	var codeInput = []int{}

	for _, i := range codeAsStrings {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		codeInput = append(codeInput, j)
	}

	codeAtZeroAddress := runcode(codeInput, 2, 12)
	fmt.Println(codeAtZeroAddress)
}

func getValue(allCode []int, paramValue int, isImmediateMode bool) int {
	if isImmediateMode {
		//fmt.Println("getValue: so ret'ing immed = ", paramValue)
		return paramValue
	}
	return allCode[paramValue]
}

func runcode(code []int, verb int, noun int) int {
	// I don't think day 5 puzzle was clear enough about whether or not to keep this part!
	// It needs to be removed.
	//code[1] = noun
	//code[2] = verb

	finalOutCode := ""

	codeIndex := 0
	for true {
		//fmt.Println("code input: ", code[codeIndex])
		newCodeIndex := -1

		paddedInstr := padInstruction(code[codeIndex])

		// find addressing modes
		// params 1 and 2 can be immediate or position.
		// param3 can never be immediate: it's only used for ADD and MUL, and these can only be sent to positions (memory)
		param2IsImmediate := (paddedInstr[1] == '1')
		param1IsImmediate := (paddedInstr[2] == '1')

		// pick out opcode
		opcodeStr := paddedInstr[3:]

		opcodeVal, err := strconv.Atoi(opcodeStr)
		if err != nil {
			panic(err)
		}
		//fmt.Println("Found opcode as int: ", opcodeVal)
		opcode := Opcode(opcodeVal)

		if opcode == END {
			fmt.Println("Halted at code: ", code)
			fmt.Println("Part 1 answer: ", finalOutCode)

			return code[0]
		}

		switch opcode {
		case ADD:
			val1 := getValue(code, code[codeIndex + 1], param1IsImmediate)
			val2 := getValue(code, code[codeIndex + 2], param2IsImmediate)
			dest := code[codeIndex + 3]

			code[dest] = val1 + val2
		case MULT:
			val1 := getValue(code, code[codeIndex + 1], param1IsImmediate)
			val2 := getValue(code, code[codeIndex + 2], param2IsImmediate)
			dest := code[codeIndex + 3]

			code[dest] = val1 * val2
			//fmt.Println("    code is now: ", code)
		case INP:
			userVal := getIntFromUser()
			dest := code[codeIndex + 1]
			code[dest] = userVal

		case OUT:
			finalOutCode := getValue(code, code[codeIndex + 1], param1IsImmediate)
			fmt.Println("------------------- OUT: ", finalOutCode)

		case JIT:
			val1 := getValue(code, code[codeIndex + 1], param1IsImmediate)
			val2 := getValue(code, code[codeIndex + 2], param2IsImmediate)

			if val1 != 0 {
				newCodeIndex = val2
			}
		case JIF:
			val1 := getValue(code, code[codeIndex + 1], param1IsImmediate)
			val2 := getValue(code, code[codeIndex + 2], param2IsImmediate)

			if val1 == 0 {
				newCodeIndex = val2
			}
		case LT:
			val1 := getValue(code, code[codeIndex + 1], param1IsImmediate)
			val2 := getValue(code, code[codeIndex + 2], param2IsImmediate)
			val3 := code[codeIndex + 3] // we know this is a position

			code[val3] = Btoi(val1 < val2)

		case EQ:
			val1 := getValue(code, code[codeIndex + 1], param1IsImmediate)
			val2 := getValue(code, code[codeIndex + 2], param2IsImmediate)
			val3 := code[codeIndex + 3] // we know this is a position

			code[val3] = Btoi(val1 == val2)

		default:
			log.Fatal("Found unrecognized opcode: ", opcode)
			panic("Abort")
		}

		if newCodeIndex >= 0 {
			// for jump instructions
			codeIndex = newCodeIndex
		} else {
			codeIndex += instrLengths[opcodeVal]
		}
	}
	return 0
}

func getIntFromUser() int {
	var i int

	fmt.Println("Please enter a number, dear astronaut: ")
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		panic(err)
	}

	return i
}
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
