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

type dot struct {
	x int
	y int
}

type dotList []dot

func (dots dotList) Len() int {
	return len(dots)
}

func (dots dotList) Less(i, j int) bool {
	if dots[i].x < dots[j].x {
		return true
	} else if dots[i].x > dots[j].x {
		return false
	} else {
		return dots[i].y < dots[j].y
	}
}

func (dots dotList) Swap(i, j int) {
	dots[i], dots[j] = dots[j], dots[i]
}

const (
	horizontal int = iota
	vertical   int = iota
)

type fold struct {
	foldType int
	coord    int
}

func scanInput(scanner *bufio.Scanner, dots *dotList, folds *[]fold) {
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}
		tokens := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])
		*dots = append(*dots, dot{x, y})
	}

	for scanner.Scan() {
		tokens := strings.Split(strings.Fields(scanner.Text())[2], "=")

		foldType := horizontal
		if tokens[0] == "y" {
			foldType = vertical
		}
		coord, _ := strconv.Atoi(tokens[1])
		*folds = append(*folds, fold{foldType, coord})
	}
}

func applyFold(dots *dotList, f fold) {
	if f.foldType == horizontal {
		for i := range *dots {
			if (*dots)[i].x >= f.coord {
				(*dots)[i].x = 2*f.coord - (*dots)[i].x
			}
		}
	} else {
		for i := range *dots {
			if (*dots)[i].y >= f.coord {
				(*dots)[i].y = 2*f.coord - (*dots)[i].y
			}
		}
	}
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	// canVisitTwice := flag.Bool("", true, "can visit a small cave twice")
	flag.Parse()

	dots := dotList{}
	folds := []fold{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &dots, &folds)
	})

	applyFold(&dots, folds[0])
	sort.Sort(dots)

	numDistinct := 1
	for i := 1; i < len(dots); i++ {
		if dots[i] != dots[i-1] {
			numDistinct++
		}
	}
	fmt.Println(numDistinct)

	for _, fold := range folds[1:] {
		applyFold(&dots, fold)
	}

	maxX := 0
	maxY := 0
	for _, dot := range dots {
		if dot.x > maxX {
			maxX = dot.x
		}
		if dot.y > maxY {
			maxY = dot.y
		}
	}

	grid := make([][]bool, maxX+1)
	for i := range grid {
		grid[i] = make([]bool, maxY+1)
	}

	for _, dot := range dots {
		grid[dot.x][dot.y] = true
	}

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			if !grid[x][y] {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Print("\n")
	}
}
