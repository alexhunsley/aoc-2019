package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	//"strings"

	//"bufio"
	//"os"
	//"strings"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////
// machine definition

type IntMachine struct {
	code *[]int
	programCounter int
	inputs []int
	output int
	didHalt bool
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// opcodes

type Opcode struct {
	code int
	// can be used for debugging
	desc string
	// Used to increment PC to reach next instruction.
	// Note you must increment PC by paramCount + 1, to allow for opcode itself.
	paramCount int
}

const (
	ADD int = iota + 1
	MULT
	INP
	OUT
	JIT
	JIF
	LT
	EQ
	HALT = 99
)

var opcodes = []Opcode{ {ADD, "ADD", 3},
	{MULT, "MULT", 3},
	{INP, "INP", 1},
	{OUT, "OUT", 1},
	{JIT, "JIT", 2},
	{JIF, "JIF", 2},
	{LT, "LT", 3},
	{EQ, "EQ", 3},
	{HALT, "HALT", 0},
}

func (Opcode) forCode(code int) Opcode {
	if code == 99 {
		return opcodes[8]
	}
	return opcodes[code - 1]
}

func padInstruction(instrCode int) string {
	return fmt.Sprintf("%05d", instrCode)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// combinatorics

// generateCombinations generates all combinations of the
// characters in the symbols string.
// For example, generateCombinations("123") would return:
//     [123 132 213 231 312 321]
//
// The combinations are generated in stable prefix order, i.e.
// the front end of the generated combinations strings changes the least
// in the order of the array returned, and the backend has the most churn.
func generateCombinations(symbols string) []string {
	resultsAccum := []string{}
	partialStr := ""

	generateCombos(&resultsAccum, partialStr, symbols)
	return resultsAccum
}

// generateCombos is a recursive internal function used by generateCombinations.
// For all your combo generating needs, please use generateCombinations.
func generateCombos(resultsAccum *[]string, partialStr string, remainingSymbols string) {
	if len(remainingSymbols) == 1 {
		partialStr += remainingSymbols
		*resultsAccum = append(*resultsAccum, partialStr)
		return
	}
	loopMax := len(remainingSymbols)
	for symbolIdx := 0; symbolIdx < loopMax; symbolIdx++ {
		//fmt.Println("SymbolIdx loop: remaining syms, i = ", remainingSymbols, symbolIdx)
		symbolToAdd := string(remainingSymbols[symbolIdx])
		partialStrAppend := partialStr
		partialStrAppend += symbolToAdd
		generateCombos(resultsAccum, partialStrAppend, strings.Replace(remainingSymbols, symbolToAdd, "", 1))
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// amplifier chain and intcode machine

func runAmplifierChain(inputPhases string, code []int) int {
	isFirstChain := true

	var machines []IntMachine

	for i := 0; i < 5; i++ {
		machineCopy := make([]int, len(code))
		copy(machineCopy, code)

		machine := IntMachine{
			code:           &machineCopy,
			programCounter: 0,
		}

		machines = append(machines, machine)
	}

	outputSignal := 0

	for true {
		for index, phase := range inputPhases {
			if isFirstChain {
				// only provide the phase in the first run of the chain
				machines[index].inputs = []int{int(phase - '0'), outputSignal}
			} else {
				machines[index].inputs = []int{outputSignal}
			}

			runcode(&(machines[index]))

			outputSignal = machines[index].output

			if machines[index].didHalt {
				return machines[4].output
			}
		}
		isFirstChain = false
	}
	// should never reach here
	return -1
}

func getValue(allCode *[]int, paramValue int, isImmediateMode bool) int {
	if isImmediateMode {
		//fmt.Println("getValue: so ret'ing immed = ", paramValue)
		return paramValue
	}
	return (*allCode)[paramValue]
}

func runcode(machine *IntMachine) {
	// IntMachine contains no heavy data itself, only small data and pointers to potentially data,
	// so this won't incurring a big copy

	//fmt.Println("runcode: given this machine: len instructions, in, out, pc: ", len(*(machine.code)), machine.inputs, machine.output, machine.programCounter)
	for true {
		//fmt.Println("code input, first codes: ", code[codeIndex], code[:10])
		newCodeIndex := -1

		codeArrayPtr := machine.code

		instructionAsStr := (*codeArrayPtr)[machine.programCounter]
		paddedInstr := padInstruction(instructionAsStr)

		// find addressing modes
		// params 1 and 2 can be immediate or position.
		// param3 can never be immediate: it's only used for ADD and MUL, and these can only be sent to positions (memory)
		param2IsImmediate := paddedInstr[1] == '1'
		param1IsImmediate := paddedInstr[2] == '1'

		// pick out opcode
		opcodeStr := paddedInstr[3:]

		opcodeVal, err := strconv.Atoi(opcodeStr)

		opcode := Opcode{}.forCode(opcodeVal)

		if err != nil {
			panic(err)
		}

		switch opcode.code {
		case ADD:
			val1 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 1], param1IsImmediate)
			val2 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 2], param2IsImmediate)
			dest := (*codeArrayPtr)[machine.programCounter + 3]

			(*codeArrayPtr)[dest] = val1 + val2
		case MULT:
			val1 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 1], param1IsImmediate)
			val2 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 2], param2IsImmediate)
			dest := (*codeArrayPtr)[machine.programCounter + 3]

			(*codeArrayPtr)[dest] = val1 * val2
		case INP:
			if len(machine.inputs) == 0 {
				panic("tried to pop an input val but inputs stack is empty")
			}
			poppedInputVal := machine.inputs[0]
			machine.inputs = machine.inputs[1:]

			dest := (*codeArrayPtr)[machine.programCounter + 1]
			//fmt.Println("writing to dest ", dest)
			(*codeArrayPtr)[dest] = poppedInputVal

		case OUT:
			machine.output = getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 1], param1IsImmediate)

			machine.programCounter += opcode.paramCount + 1
			machine.didHalt = false

			return

		case JIT:
			val1 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 1], param1IsImmediate)
			val2 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 2], param2IsImmediate)

			if val1 != 0 {
				newCodeIndex = val2
			}
		case JIF:
			val1 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 1], param1IsImmediate)
			val2 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 2], param2IsImmediate)

			if val1 == 0 {
				newCodeIndex = val2
			}
		case LT:
			val1 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 1], param1IsImmediate)
			val2 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 2], param2IsImmediate)
			// we know this is a position
			val3 := (*codeArrayPtr)[machine.programCounter + 3]

			(*codeArrayPtr)[val3] = Btoi(val1 < val2)

		case EQ:
			val1 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 1], param1IsImmediate)
			val2 := getValue(codeArrayPtr, (*codeArrayPtr)[machine.programCounter + 2], param2IsImmediate)
			// we know this is a position
			val3 := (*codeArrayPtr)[machine.programCounter + 3]

			(*codeArrayPtr)[val3] = Btoi(val1 == val2)

		case HALT:
			machine.programCounter += opcode.paramCount + 1

			machine.didHalt = true
		default:
			log.Fatal("Found unrecognized opcode: ", opcode)
			panic("Abort")
		}

		if newCodeIndex >= 0 {
			// for jump instructions
			machine.programCounter = newCodeIndex
		} else {
			machine.programCounter += opcode.paramCount + 1
		}

		if machine.didHalt {
			return
		}
	}
	return
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// helper

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// main

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

	var code = []int{}

	for _, i := range codeAsStrings {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		code = append(code, j)
	}

	allPhaseCombos := generateCombinations("56789")

	maxPower := 0

	for _, phases := range allPhaseCombos {
		power := runAmplifierChain(phases, code)
		if power > maxPower {
			maxPower = power
		}
	}

	fmt.Println("Part 2 solution (max power): ", maxPower)
}
