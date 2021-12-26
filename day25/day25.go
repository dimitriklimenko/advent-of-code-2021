package main

import (
	"bufio"
	"flag"
	"fmt"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type state [][]byte

func scanInput(scanner *bufio.Scanner, result *state) {
	for scanner.Scan() {
		*result = append(*result, []byte(scanner.Text()))
	}
}

func printState(s state) {
	for _, row := range s {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func step(s state) bool {
	movedEast := stepEast(s)
	movedSouth := stepSouth(s)
	return movedEast || movedSouth
}

func stepEast(s state) bool {
	moved := false
	numRows := len(s)
	numCols := len(s[0])
	for i := 0; i < numRows; i++ {
		movers := []int{}
		for j := 0; j < numCols; j++ {
			if s[i][j] == '>' && s[i][(j+1)%numCols] == '.' {
				movers = append(movers, j)
				moved = true
			}
		}
		for _, j := range movers {
			s[i][(j+1)%numCols] = '>'
			s[i][j] = '.'
		}
	}
	return moved
}

func stepSouth(s state) bool {
	moved := false
	numRows := len(s)
	numCols := len(s[0])
	for j := 0; j < numCols; j++ {
		movers := []int{}
		for i := 0; i < numRows; i++ {
			if s[i][j] == 'v' && s[(i+1)%numRows][j] == '.' {
				movers = append(movers, i)
				moved = true
			}
		}
		for _, i := range movers {
			s[(i+1)%numRows][j] = 'v'
			s[i][j] = '.'
		}
	}
	return moved
}

func main() {
	filePath := flag.String("f", "", "file to read from / write to")
	flag.Parse()

	currentState := state{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &currentState)
	})

	numSteps := 0
	for {
		// printState(currentState)
		moved := step(currentState)
		numSteps++
		if !moved {
			break
		}
	}
	// printState(currentState)
	fmt.Println(numSteps)
}
