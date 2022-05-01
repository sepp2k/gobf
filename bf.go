package main

import (
	"fmt"
	"os"
)

type OpCode int

const (
	Add OpCode = iota
	Move
	Out
	In
	JumpIfZero
	JumpIfNonZero
)

type Instruction struct {
	opCode  OpCode
	operand int
}

func incrementInstruction(instructions []Instruction, opCode OpCode, increment int) []Instruction {
	last := len(instructions) - 1
	if last >= 0 && instructions[last].opCode == opCode {
		instructions[last].operand += increment
		return instructions
	} else {
		return append(instructions, Instruction{opCode, increment})
	}
}

func parse(source []byte) []Instruction {
	instructions := make([]Instruction, 0, 1024)
	loopStack := make([]int, 0, 16)
	for _, c := range source {
		switch c {
		case '[':
			loopStack = append(loopStack, len(instructions))
			instructions = append(instructions, Instruction{JumpIfZero, -1})
		case ']':
			n := len(loopStack) - 1
			opening := loopStack[n]
			instructions = append(instructions, Instruction{JumpIfNonZero, opening + 1})
			instructions[opening].operand = len(instructions)
			loopStack = loopStack[:n]
		case '+':
			instructions = incrementInstruction(instructions, Add, 1)
		case '-':
			instructions = incrementInstruction(instructions, Add, -1)
		case '<':
			instructions = incrementInstruction(instructions, Move, -1)
		case '>':
			instructions = incrementInstruction(instructions, Move, 1)
		case '.':
			instructions = append(instructions, Instruction{Out, -1})
		case ',':
			instructions = append(instructions, Instruction{In, -1})
		}
	}
	return instructions
}

func readTape(tape []byte, index int) byte {
	if index < len(tape) {
		return tape[index]
	} else {
		return 0
	}
}

func writeTape(tape *[]byte, index int, value byte) {
	n := len(*tape)
	if index < n {
		(*tape)[index] = value
	} else {
		*tape = append(*tape, make([]byte, index-n+1)...)
	}
}

func exec(instructions []Instruction) {
	tape := make([]byte, 1024)
	tapePointer := 0
	for pc := 0; pc < len(instructions); pc++ {
		switch instructions[pc].opCode {
		case Out:
			fmt.Printf("%c", readTape(tape, tapePointer))
		case In:
			writeTape(&tape, tapePointer, 0)
			_, err := fmt.Scanf("%c", &tape[tapePointer])
			if err != nil {
				panic(err)
			}
		case Add:
			writeTape(&tape, tapePointer, readTape(tape, tapePointer)+byte(instructions[pc].operand))
		case Move:
			tapePointer += instructions[pc].operand
		case JumpIfZero:
			if tape[tapePointer] == 0 {
				pc = instructions[pc].operand - 1 // Subtract 1 to account for loop increment
			}
		case JumpIfNonZero:
			if tape[tapePointer] != 0 {
				pc = instructions[pc].operand - 1 // Subtract 1 to account for loop increment
			}
		}
	}
}

func main() {
	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	instructions := parse(source)
	exec(instructions)
}
