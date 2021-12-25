package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type Point struct {
	x int
	y int
}

type Line struct {
	p1 Point
	p2 Point
}

type Grid [1000][1000]int

func parseLine(text string) Line {
	parts := strings.Split(text, "->")
	return Line{
		parsePoint(parts[0]),
		parsePoint(parts[1]),
	}
}

func parsePoint(text string) Point {
	parts := strings.Split(text, ",")
	x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return Point{
		x,
		y,
	}
}

func scanInput(scanner *bufio.Scanner, lines *[]Line) {
	for scanner.Scan() {
		*lines = append(*lines, parseLine(scanner.Text()))
	}
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	countDiagonals := flag.Bool("last", true, "set to count diagonals")
	flag.Parse()

	lines := []Line{}

	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &lines)
	})

	grid := Grid{}

	for _, line := range lines {
		var lo int
		var hi int
		if line.p1.x == line.p2.x {
			if line.p1.y <= line.p2.y {
				lo = line.p1.y
				hi = line.p2.y
			} else {
				lo = line.p2.y
				hi = line.p1.y
			}
			for y := lo; y <= hi; y++ {
				grid[line.p1.x][y]++
			}
		} else if line.p1.y == line.p2.y {
			if line.p1.x <= line.p2.x {
				lo = line.p1.x
				hi = line.p2.x
			} else {
				lo = line.p2.x
				hi = line.p1.x
			}
			for x := lo; x <= hi; x++ {
				grid[x][line.p1.y]++
			}
		} else if *countDiagonals {
			dx := 1
			numSteps := line.p2.x - line.p1.x
			if line.p1.x > line.p2.x {
				dx = -1
				numSteps = line.p1.x - line.p2.x
			}
			dy := 1
			if line.p1.y > line.p2.y {
				dy = -1
			}

			x := line.p1.x
			y := line.p1.y
			for i := 0; i <= numSteps; i++ {
				grid[x][y]++
				x += dx
				y += dy
			}

		}
	}

	numOverlaps := 0
	for _, row := range grid {
		for _, count := range row {
			if count >= 2 {
				numOverlaps++
			}
		}

	}
	fmt.Println(numOverlaps)
}
