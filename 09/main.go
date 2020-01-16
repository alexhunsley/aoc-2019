package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////
// machine definition

type IntMachine struct {
	code *[]int
	programCounter int
	inputs []int
	output int
	// map from address to value
	sparseMemory *map[int]int
	relativeBase int
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
	ADDR_MODE_POSITION uint8 = iota
	ADDR_MODE_IMMEDIATE
	ADDR_MODE_RELATIVE
)

const (
	ADD int = iota + 1
	MULT
	INP
	OUT
	JIT
	JIF
	LT
	EQ
	ARB
	HALT = 99
)

var opcodes = []Opcode{
	{HALT, "HALT", 0},
	{ADD, "ADD", 3},
	{MULT, "MULT", 3},
	{INP, "INP", 1},
	{OUT, "OUT", 1},
	{JIT, "JIT", 2},
	{JIF, "JIF", 2},
	{LT, "LT", 3},
	{EQ, "EQ", 3},
	{ARB, "ARB", 1},
}

func (Opcode) forCode(code int) Opcode {
	if code == 99 {
		return opcodes[0]
	}
	return opcodes[code]
}

func padInstruction(instrCode int) string {
	return fmt.Sprintf("%05d", instrCode)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// intcode machine

func getValue(machine *IntMachine, paramValue int, mode uint8) int {
	//fmt.Println("get value, mode = ", mode)
	if mode == ADDR_MODE_IMMEDIATE {
		//fmt.Println("getValue: so ret'ing immed = ", paramValue)
		return paramValue
	}

	// it's position or relative address.
	if mode == ADDR_MODE_RELATIVE {
		// made relative mode adjustment to address
		paramValue += machine.relativeBase
	}

	if paramValue < 0 {
		panic("Tried to access memory at negative address")
	}

	if paramValue >= len(*machine.code) {
		//fmt.Println("fetching sparse value at addr, with sparse contents =  ", paramValue, machine.sparseMemory)
		// it's a sparse memory value
		return (*machine.sparseMemory)[paramValue]
	}
	return (*machine.code)[paramValue]
}

func setValue(machine *IntMachine, address int, value int, mode uint8) {
	if mode == ADDR_MODE_IMMEDIATE {
		panic("Attempted to store a value to an immediate mode value (rather than position or relative)")
	}
	// it's position or relative address.
	if mode == ADDR_MODE_RELATIVE {
		// made relative mode adjustment to address
		address += machine.relativeBase
	}

	if address < 0 {
		panic("Tried to set memory at negative address")
	}

	// check if goes off end of the memory
	if address >= len(*machine.code) {
		//fmt.Println("storing sparse,  value,  at addr = ", machine.sparseMemory, value, address)
		(*machine.sparseMemory)[address] = value
		return
	}
	(*machine.code)[address] = value
}

func formatInstrWithParams(machine *IntMachine, values []int, paramModes []uint8) string {
	var str strings.Builder

	// prefixes to params: position, #immediate, +relative
	addressModePrefixes := []string{"", "#", "~"}

	for i := 0; i < len(values); i++ {
		str.WriteString(addressModePrefixes[paramModes[i]])
		str.WriteString(strconv.Itoa(values[i]))
		str.WriteString(", ")
	}

	// cut off the last ", " part
	resultStr := str.String()

	//fmt.Println("before cut: ", resultStr)
	if len(values) > 0 {
		resultStr = resultStr[:len(resultStr) - 2]
	}
	//fmt.Println("AFTER cut: ", resultStr)
	resultStr += "   (arb = "
	resultStr += strconv.Itoa(machine.relativeBase)
	resultStr += ", pModes = "
	resultStr += fmt.Sprintf("%s", paramModes)
	resultStr += ")"

	return resultStr
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
		param3Mode := paddedInstr[0] - '0'
		param2Mode := paddedInstr[1] - '0'
		param1Mode := paddedInstr[2] - '0'

		// pick out opcode
		opcodeStr := paddedInstr[3:]

		opcodeVal, err := strconv.Atoi(opcodeStr)

		opcode := Opcode{}.forCode(opcodeVal)

		//fmt.Println(" getting op from addr: ", machine.programCounter, " opcode: -- ", opcode.desc, param1Mode, param2Mode, "entire code = ", machine.code)
		if err != nil {
			panic(err)
		}

		switch opcode.code {
		case ADD:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)
			dest := (*codeArrayPtr)[machine.programCounter + 3]

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("ADD ", formatInstrWithParams(machine, []int{val1, val2, dest}, []uint8{param1Mode, param2Mode, param3Mode}))

			//fmt.Printf("====== addr %d  ADD %d, %d, %d (%d %d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], (*codeArrayPtr)[machine.programCounter + 2], dest, param1Mode, param2Mode)

			setValue(machine, dest, val1 + val2, param3Mode)
		case MULT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)
			dest := (*codeArrayPtr)[machine.programCounter + 3]

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("MUL ", formatInstrWithParams(machine, []int{val1, val2, dest}, []uint8{param1Mode, param2Mode, param3Mode}))

			//fmt.Printf("====== addr %d  MULT %d, %d, %d (%d %d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], (*codeArrayPtr)[machine.programCounter + 2], dest, param1Mode, param2Mode)

			setValue(machine, dest, val1 * val2, param3Mode)
		case INP:
			// Have to be really careful with this one.
			// We have to adjust the relative base 'directly' in here; calling getValue does the wrong thing.

			if len(machine.inputs) == 0 {
				panic("tried to pop an input val but inputs stack is empty")
			}
			poppedInputVal := machine.inputs[0]
			machine.inputs = machine.inputs[1:]


			//fmt.Printf("====== addr %d  INP %d (%d)   ((rel base = %d))\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], param1Mode, machine.relativeBase)

			dest := (*codeArrayPtr)[machine.programCounter + 1]

			dest += machine.relativeBase

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("INP ", formatInstrWithParams(machine, []int{dest}, []uint8{ADDR_MODE_RELATIVE}))

			setValue(machine, dest, poppedInputVal, ADDR_MODE_POSITION)

			//fmt.Println("writing to dest ", dest)

		case OUT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			machine.output = val1

			//fmt.Printf("====== addr %d  OUT %d (%d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("OUTPUT ", formatInstrWithParams(machine, []int{val1}, []uint8{param1Mode}))

		case JIT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)

			//fmt.Printf("====== addr %d  JIT %d %d (%d %d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], (*codeArrayPtr)[machine.programCounter + 2], param1Mode, param2Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("JIT ", formatInstrWithParams(machine, []int{val1, val2}, []uint8{param1Mode, param2Mode}))

			if val1 != 0 {
				log.Debug("-------------------------------------------------------------------------------------------")
				newCodeIndex = val2
			}
		case JIF:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)

			//fmt.Printf("====== addr %d  JIF %d %d (%d %d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], (*codeArrayPtr)[machine.programCounter + 2], param1Mode, param2Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("JIF ", formatInstrWithParams(machine, []int{val1, val2}, []uint8{param1Mode, param2Mode}))

			if val1 == 0 {
				log.Debug("-------------------------------------------------------------------------------------------")
				newCodeIndex = val2
			}
		case LT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)
			// we know this is a position
			val3 := (*codeArrayPtr)[machine.programCounter + 3]

			//fmt.Printf("====== addr %d  LT %d %d %d (%d %d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], (*codeArrayPtr)[machine.programCounter + 2], (*codeArrayPtr)[machine.programCounter + 3], param1Mode, param2Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("LT ", formatInstrWithParams(machine, []int{val1, val2, val3}, []uint8{param1Mode, param2Mode, param3Mode}))

			setValue(machine, val3, Btoi(val1 < val2), param3Mode)

		case EQ:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)
			// we know this is a position
			val3 := (*codeArrayPtr)[machine.programCounter + 3]

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("EQ ", formatInstrWithParams(machine, []int{val1, val2, val3}, []uint8{param1Mode, param2Mode, param3Mode}))

			//fmt.Printf("====== addr %d  EQ %d %d %d (%d %d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], (*codeArrayPtr)[machine.programCounter + 2], (*codeArrayPtr)[machine.programCounter + 3], param1Mode, param2Mode)
			setValue(machine, val3, Btoi(val1 == val2), param3Mode)

		case ARB:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)

			//fmt.Printf("====== addr %d  OUT %d (%d)\n", machine.programCounter, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("ARB ", formatInstrWithParams(machine, []int{val1}, []uint8{param1Mode}), "  and raw param1 before RBO is ", (*codeArrayPtr)[machine.programCounter + 1])

			machine.relativeBase += val1

		case HALT:
			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Debug("HALT")

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
	//log.SetLevel(log.ErrorLevel)
	log.SetLevel(log.DebugLevel)
	logFile, logErr := os.OpenFile("go.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if logErr == nil {
		log.SetOutput(logFile)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.close()
	
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

	part1Code := make([]int, len(code))
	copy(part1Code, code)

	machine := IntMachine{
		code:           &part1Code,
		programCounter: 0,
		inputs:	[]int{1},
		sparseMemory: &(map[int]int{}),
	}

	runcode(&machine)

	part1Soln := machine.output

	part2Code := make([]int, len(code))
	copy(part2Code, code)

	machine2 := IntMachine{
		code:           &part2Code,
		programCounter: 0,
		inputs:	[]int{2},
		sparseMemory: &(map[int]int{}),
	}

	runcode(&machine2)

	log.Info("Part 1 solution: ", part1Soln)
	log.Info("Part 2 solution: ", machine2.output)
}
