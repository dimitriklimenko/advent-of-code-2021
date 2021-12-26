package main

import (
	"bufio"
	"flag"
	"fmt"
	"regexp"
	"strconv"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

var lineRegex = regexp.MustCompile("(on|off) x=(-?[0-9]+)..(-?[0-9]+),y=(-?[0-9]+)..(-?[0-9]+),z=(-?[0-9]+)..(-?[0-9]+)$")

type coordRange struct {
	min int
	max int
}

type instruction struct {
	turnOn bool
	xRange coordRange
	yRange coordRange
	zRange coordRange
}

type limitedGrid [101][101][101]bool

func parseRange(parts [][]byte) coordRange {
	rng := coordRange{}
	rng.min, _ = strconv.Atoi(string(parts[0]))
	rng.max, _ = strconv.Atoi(string(parts[1]))
	return rng
}

func scanInput(scanner *bufio.Scanner, steps *[]instruction) {
	for scanner.Scan() {
		text := scanner.Text()
		bytes := []byte(text)
		matches := lineRegex.FindSubmatch(bytes)
		if len(matches) != 8 {
			panic("Invalid input!?")
		}
		step := instruction{
			turnOn: string(matches[1]) == "on",
			xRange: parseRange(matches[2:4]),
			yRange: parseRange(matches[4:6]),
			zRange: parseRange(matches[6:8]),
		}
		*steps = append(*steps, step)
	}
}

func boundLow(value int) int {
	if value < -50 {
		value = -50
	}
	return value
}

func boundHigh(value int) int {
	if value > 50 {
		value = 50
	}
	return value
}

func processStep(grid *limitedGrid, step instruction) {
	for i := boundLow(step.xRange.min) + 50; i <= boundHigh(step.xRange.max)+50; i++ {
		for j := boundLow(step.yRange.min) + 50; j <= boundHigh(step.yRange.max)+50; j++ {
			for k := boundLow(step.zRange.min) + 50; k <= boundHigh(step.zRange.max)+50; k++ {
				(*grid)[i][j][k] = step.turnOn
			}
		}
	}
}

func processSteps(grid *limitedGrid, steps []instruction) {
	for _, step := range steps {
		processStep(grid, step)
		// fmt.Println(countCubesOn(*grid))
	}
}

func countCubesOn(grid limitedGrid) int {
	count := 0
	for i := 0; i < 101; i++ {
		for j := 0; j < 101; j++ {
			for k := 0; k < 101; k++ {
				if grid[i][j][k] {
					count++
				}
			}
		}
	}
	return count
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	// isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	steps := []instruction{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &steps)
	})
	fmt.Println()

	grid := limitedGrid{}
	processSteps(&grid, steps)
	count := countCubesOn(grid)
	fmt.Println(count)
}
