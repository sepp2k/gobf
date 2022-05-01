package main

import (
	"fmt"
	"os"
)

func main() {
	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	loopStack := make([]int, 0, 16)
	loopJumps := make(map[int]int)
	for i, c := range source {
		if c == '[' {
			loopStack = append(loopStack, i)
		} else if c == ']' {
			n := len(loopStack) - 1
			opening := loopStack[n]
			loopJumps[i] = opening
			loopJumps[opening] = i
			loopStack = loopStack[:n]
		}
	}

	tape := make([]byte, 1024)
	tapePointer := 0
	for pc := 0; pc < len(source); pc++ {
		switch source[pc] {
		case '.':
			fmt.Printf("%c", tape[tapePointer])
		case ',':
			_, err := fmt.Scanf("%c", &tape[tapePointer])
			if err != nil {
				panic(err)
			}
		case '+':
			tape[tapePointer]++
		case '-':
			tape[tapePointer]--
		case '>':
			tapePointer++
			if tapePointer >= len(tape) {
				tape = append(tape, 0)
			}
		case '<':
			tapePointer--
		case '[':
			if tape[tapePointer] == 0 {
				pc = loopJumps[pc]
			}
		case ']':
			if tape[tapePointer] != 0 {
				pc = loopJumps[pc]
			}
		}
	}
}
