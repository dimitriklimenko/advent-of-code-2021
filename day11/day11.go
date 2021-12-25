package main

import (
	"bufio"
	"flag"
	"fmt"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type octopusGrid [10][10]int
type flashGrid [10][10]bool

type location struct {
	i int
	j int
}

func getNeighbours(loc location) []location {
	nbs := []location{}
	for _, di := range []int{-1, 0, +1} {
		if loc.i+di < 0 || loc.i+di >= 10 {
			continue
		}
		for _, dj := range []int{-1, 0, +1} {
			if loc.j+dj < 0 || loc.j+dj >= 10 {
				continue
			}
			if di == 0 && dj == 0 {
				continue
			}
			nbs = append(nbs, location{loc.i + di, loc.j + dj})
		}
	}
	return nbs
}

func scanInput(scanner *bufio.Scanner, grid *octopusGrid) {
	i := 0
	for scanner.Scan() {
		for j, c := range scanner.Text() {
			grid[i][j] = int(c - '0')
		}
		i += 1
	}
}

func stepGrid(grid *octopusGrid) int {
	numFlashes := 0
	flashes := flashGrid{}
	flashQueue := []location{}
	for i, row := range grid {
		for j, val := range row {
			grid[i][j] = val + 1
			if val+1 > 9 {
				flashQueue = append(flashQueue, location{i, j})
			}
		}
	}

	for queueIndex := 0; queueIndex < len(flashQueue); queueIndex++ {
		loc := flashQueue[queueIndex]
		if flashes[loc.i][loc.j] {
			continue
		}
		flashes[loc.i][loc.j] = true
		numFlashes++

		for _, nb := range getNeighbours(loc) {
			grid[nb.i][nb.j] += 1
			if grid[nb.i][nb.j] > 9 && !flashes[nb.i][nb.j] {
				flashQueue = append(flashQueue, nb)
			}
		}
	}

	for i, row := range grid {
		for j, val := range row {
			if val > 9 {
				grid[i][j] = 0
			}
		}
	}

	return numFlashes
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	grid := octopusGrid{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &grid)
	})

	numFlashes := 0
	if !*isPart2 {
		for step := 0; step < 100; step++ {
			numFlashes += stepGrid(&grid)
		}
		fmt.Println(numFlashes)
	} else {
		for step := 1; ; step++ {
			numFlashes = stepGrid(&grid)
			if numFlashes == 100 {
				fmt.Println(step)
				break
			}
		}
	}
}
