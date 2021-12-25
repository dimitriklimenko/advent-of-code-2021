package main

import (
	"bufio"
	"flag"
	"fmt"
	"sort"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type heightMap [][]int
type location struct {
	i int
	j int
}

func scanInput(scanner *bufio.Scanner, heights *heightMap) {
	for scanner.Scan() {
		row := []int{}
		for _, c := range scanner.Text() {
			row = append(row, int(c-'0'))
		}
		*heights = append(*heights, row)
	}
}

func getNeighbours(loc location, numRows int, numCols int) []location {
	nbs := []location{}
	if loc.i > 0 {
		nbs = append(nbs, location{loc.i - 1, loc.j})
	}
	if loc.i < numRows-1 {
		nbs = append(nbs, location{loc.i + 1, loc.j})
	}
	if loc.j > 0 {
		nbs = append(nbs, location{loc.i, loc.j - 1})
	}
	if loc.j < numCols-1 {
		nbs = append(nbs, location{loc.i, loc.j + 1})
	}
	return nbs
}

func findLows(heights heightMap) []location {
	locations := []location{}
	numRows := len(heights)
	numCols := len(heights[0])
	for i, row := range heights {
		for j, val := range row {
			isLowPoint := true
			for _, nb := range getNeighbours(location{i, j}, numRows, numCols) {
				if val >= heights[nb.i][nb.j] {
					isLowPoint = false
					break
				}
			}
			if isLowPoint {
				locations = append(locations, location{i, j})
			}
		}
	}
	return locations
}

func getBasinSize(heights heightMap, lowLocation location) int {
	numRows := len(heights)
	numCols := len(heights[0])

	basinMap := make([][]bool, numRows)
	for i := range basinMap {
		basinMap[i] = make([]bool, numCols)
	}

	basinSize := 0
	queue := make(chan location, numRows*numCols)
	basinMap[lowLocation.i][lowLocation.j] = true
	queue <- lowLocation

	notEmpty := true
	for notEmpty {
		select {
		case loc := <-queue:
			basinSize++
			for _, nb := range getNeighbours(loc, numRows, numCols) {
				if heights[nb.i][nb.j] < 9 && !basinMap[nb.i][nb.j] {
					basinMap[nb.i][nb.j] = true
					queue <- nb
				}
			}
		default:
			notEmpty = false
		}
	}
	return basinSize
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	var heights heightMap
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &heights)
	})
	lows := findLows(heights)

	if !*isPart2 {
		totalRisk := 0
		for _, loc := range lows {
			totalRisk += heights[loc.i][loc.j] + 1
		}
		fmt.Println(totalRisk)
	} else {
		basinSizes := []int{}
		for _, loc := range lows {
			sz := getBasinSize(heights, loc)
			basinSizes = append(basinSizes, sz)
		}
		sort.Ints(basinSizes)
		nb := len(lows)
		sizeProduct := basinSizes[nb-1] * basinSizes[nb-2] * basinSizes[nb-3]
		fmt.Println(sizeProduct)
	}
}
