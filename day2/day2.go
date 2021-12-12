package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

func scanInstructions(scanner *bufio.Scanner, dirs *[]string, deltas *[]int) {
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		*dirs = append(*dirs, scanner.Text())
		if !scanner.Scan() {
			panic("malformed input!")
		}
		delta, _ := strconv.Atoi(scanner.Text())
		*deltas = append(*deltas, delta)
	}
}

func calcFinalProduct(dirs []string, deltas []int) int {
	horizontalPos := 0
	depth := 0
	for i := 0; i < len(dirs); i++ {
		if dirs[i] == "forward" {
			horizontalPos += deltas[i]
		} else if dirs[i] == "down" {
			depth += deltas[i]
		} else if dirs[i] == "up" {
			depth -= deltas[i]
		} else {
			panic(fmt.Sprintf("Invalid direction: %s", dirs[i]))
		}
	}
	return horizontalPos * depth
}

func calcFinalProductPartTwo(dirs []string, deltas []int) int {
	aim := 0
	horizontalPos := 0
	depth := 0
	for i := 0; i < len(dirs); i++ {
		if dirs[i] == "forward" {
			horizontalPos += deltas[i]
			depth += aim * deltas[i]
		} else if dirs[i] == "down" {
			aim += deltas[i]
		} else if dirs[i] == "up" {
			aim -= deltas[i]
		} else {
			panic(fmt.Sprintf("Invalid direction: %s", dirs[i]))
		}
	}
	return horizontalPos * depth
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	hasAim := flag.Bool("aim", true, "set for aim-based navigation")
	flag.Parse()

	dirs := []string{}
	deltas := []int{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInstructions(scanner, &dirs, &deltas)
	})

	var finalProduct int
	if *hasAim {
		finalProduct = calcFinalProductPartTwo(dirs, deltas)
	} else {
		finalProduct = calcFinalProduct(dirs, deltas)
	}
	fmt.Println(finalProduct)
}
