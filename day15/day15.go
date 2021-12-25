package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

// copied from https://pkg.go.dev/container/heap#example-package-PriorityQueue
// An Item is something we manage in a priority queue.
type Item struct {
	value location // The value of the item; arbitrary.
	risk  int      // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Lowest risk first
	return pq[i].risk < pq[j].risk
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type location struct {
	i int
	j int
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

func scanInput(scanner *bufio.Scanner, g *[][]int) {
	i := 0
	for scanner.Scan() {
		row := []int{}
		for _, c := range scanner.Text() {
			row = append(row, int(c-'0'))
		}
		*g = append(*g, row)
		i += 1
	}
}

func findLowestRisk(grid [][]int) int {
	numRows := len(grid)
	numCols := len(grid[0])
	start := location{0, 0}
	end := location{numRows - 1, numCols - 1}

	visited := make([][]bool, numRows)
	for i := range visited {
		visited[i] = make([]bool, numCols)
	}

	pq := PriorityQueue{}
	heap.Push(&pq, &Item{start, 0, -1})
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if item.value == end {
			return item.risk
		}
		if visited[item.value.i][item.value.j] {
			continue
		}
		visited[item.value.i][item.value.j] = true
		for _, nb := range getNeighbours(item.value, numRows, numCols) {
			heap.Push(&pq, &Item{nb, item.risk + grid[nb.i][nb.j], -1})
		}
	}
	return -1
}

func makeLargerGrid(grid [][]int) [][]int {
	numRows := len(grid)
	numCols := len(grid[0])
	largerGrid := make([][]int, numRows*5)
	for i := range largerGrid {
		largerGrid[i] = make([]int, numCols*5)
	}

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			for ii := 0; ii < 5; ii++ {
				for jj := 0; jj < 5; jj++ {
					value := (grid[i][j] + ii + jj) % 9
					if value == 0 {
						value = 9
					}
					largerGrid[ii*numRows+i][jj*numCols+j] = value
				}
			}
		}
	}
	return largerGrid
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	grid := [][]int{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &grid)
	})

	if *isPart2 {
		grid = makeLargerGrid(grid)
	}

	lowestRisk := findLowestRisk(grid)
	fmt.Println(lowestRisk)
}
