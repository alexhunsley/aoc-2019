package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type CallbackForGetInput func() int
type CallbackForOutput func(int)

type IntMachine struct {
	code *[]int
	programCounter int
	getInputCallback CallbackForGetInput
	output int
	// map from address to value
	sparseMemory *map[int]int
	relativeBase int
}

func (m IntMachine) poke(addr int, val int) {
	(*m.code)[addr] = val
}

func (m IntMachine) peek(addr int) int {
	return (*m.code)[addr]
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
	NOP = 98
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
	{NOP, "NOP", 0},
}

func (Opcode) forCode(code int) Opcode {
	if code == 99 {
		return opcodes[0]
	} else if code == 98 {
		return opcodes[10]
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

func runcode(machine *IntMachine, getInputCallback CallbackForGetInput, sendOutputCallback CallbackForOutput) {

runcodeLoop:
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
			}).Trace("ADD ", formatInstrWithParams(machine, []int{val1, val2, dest}, []uint8{param1Mode, param2Mode, param3Mode}))

			setValue(machine, dest, val1 + val2, param3Mode)
		case MULT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)
			dest := (*codeArrayPtr)[machine.programCounter + 3]

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("MUL ", formatInstrWithParams(machine, []int{val1, val2, dest}, []uint8{param1Mode, param2Mode, param3Mode}))

			setValue(machine, dest, val1 * val2, param3Mode)
		case INP:
			inputVal := getInputCallback()

			dest := (*codeArrayPtr)[machine.programCounter + 1]
			//dest += machine.relativeBase

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("INP ", formatInstrWithParams(machine, []int{dest}, []uint8{ADDR_MODE_RELATIVE}))

			setValue(machine, dest, inputVal, ADDR_MODE_POSITION)

		case OUT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)

			sendOutputCallback(val1)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("OUTPUT ", formatInstrWithParams(machine, []int{val1}, []uint8{param1Mode}))

		case JIT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("JIT ", formatInstrWithParams(machine, []int{val1, val2}, []uint8{param1Mode, param2Mode}))

			if val1 != 0 {
				log.Trace("-------------------------------------------------------------------------------------------")
				newCodeIndex = val2
			}
		case JIF:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("JIF ", formatInstrWithParams(machine, []int{val1, val2}, []uint8{param1Mode, param2Mode}))

			if val1 == 0 {
				log.Trace("-------------------------------------------------------------------------------------------")
				newCodeIndex = val2
			}
		case LT:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)
			// we know this is a position
			val3 := (*codeArrayPtr)[machine.programCounter + 3]

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("LT ", formatInstrWithParams(machine, []int{val1, val2, val3}, []uint8{param1Mode, param2Mode, param3Mode}))

			setValue(machine, val3, Btoi(val1 < val2), param3Mode)

		case EQ:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)
			val2 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 2], param2Mode)
			// we know this is a position
			val3 := (*codeArrayPtr)[machine.programCounter + 3]

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("EQ ", formatInstrWithParams(machine, []int{val1, val2, val3}, []uint8{param1Mode, param2Mode, param3Mode}))

			setValue(machine, val3, Btoi(val1 == val2), param3Mode)

		case ARB:
			val1 := getValue(machine, (*codeArrayPtr)[machine.programCounter + 1], param1Mode)

			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("ARB ", formatInstrWithParams(machine, []int{val1}, []uint8{param1Mode}), "  and raw param1 before RBO is ", (*codeArrayPtr)[machine.programCounter + 1])

			machine.relativeBase += val1

		case NOP:
			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("NOP")

			fmt.Println("NOP reached at instruction ", machine.programCounter)
			machine.programCounter += 1

		case HALT:
			log.WithFields(log.Fields{
				"pc": machine.programCounter,
			}).Trace("HALT")

			fmt.Println("HALT reached at instruction ", machine.programCounter)
			machine.programCounter += 1

			break runcodeLoop
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
	}
	return
}
