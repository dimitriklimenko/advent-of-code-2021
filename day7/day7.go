package main

import (
	"bufio"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

func scanInput(scanner *bufio.Scanner, crabPositions *[]int) {
	scanner.Scan()
	for _, token := range strings.Split(scanner.Text(), ",") {
		pos, _ := strconv.Atoi(token)
		*crabPositions = append(*crabPositions, pos)
	}
}

func calcCostsLinear(crabPositions []int) []int {
	totalCrabs := len(crabPositions)
	moveCosts := make([]int, totalCrabs)

	for _, step := range []int{-1, 1} {
		start := 0
		end := totalCrabs
		if step == -1 {
			start = totalCrabs - 1
			end = -1
		}

		numCrabs := 0
		totalCost := 0
		lastPos := 0

		for i := start; i != end; i += step {
			pos := crabPositions[i]
			var delta int
			if step == 1 {
				delta = pos - lastPos
			} else {
				delta = lastPos - pos
			}
			lastPos = pos
			totalCost += numCrabs * delta
			numCrabs++
			moveCosts[i] += totalCost
		}
	}
	return moveCosts
}

func calcMin(values []int) int {
	min := values[0]
	for _, val := range values {
		if val < min {
			min = val
		}
	}
	return min
}

func calcMinQuadratic(crabPositions []int) int {
	totalCrabs := len(crabPositions)
	lo := crabPositions[0]
	hi := crabPositions[totalCrabs-1]

	// Guaranteed to be higher than the minimum cost...
	minCost := totalCrabs * (hi - lo) * (hi - lo + 1)

	for targetPos := lo; targetPos <= hi; targetPos++ {
		totalCost := 0
		for _, pos := range crabPositions {
			delta := pos - targetPos
			if delta < 0 {
				delta = -delta
			}
			totalCost += delta * (delta + 1) / 2
		}
		if totalCost < minCost {
			minCost = totalCost
		}
	}
	return minCost
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	hasQuadraticCost := flag.Bool("q", true, "true iff the fuel costs are quadratic")
	flag.Parse()

	crabPositions := []int{}

	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &crabPositions)
	})
	sort.Ints(crabPositions)

	var minCost int
	if *hasQuadraticCost {
		minCost = calcMinQuadratic(crabPositions)
	} else {
		minCost = calcMin(calcCostsLinear(crabPositions))
	}

	fmt.Println(minCost)
}
