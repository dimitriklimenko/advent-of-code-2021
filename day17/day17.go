package main

import (
	"bufio"
	"flag"
	"fmt"
	"regexp"
	"strconv"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

var targetRegex = regexp.MustCompile("^target area: x=(-?[0-9]+)..(-?[0-9]+), y=(-?[0-9]+)..(-?[0-9]+)$")

func scanInput(scanner *bufio.Scanner, minX *int, maxX *int, minY *int, maxY *int) {
	scanner.Scan()
	text := scanner.Text()
	bytes := []byte(text)
	matches := targetRegex.FindSubmatch(bytes)
	if len(matches) != 5 {
		panic("Invalid input!?")
	}
	*minX, _ = strconv.Atoi(string(matches[1]))
	*maxX, _ = strconv.Atoi(string(matches[2]))
	*minY, _ = strconv.Atoi(string(matches[3]))
	*maxY, _ = strconv.Atoi(string(matches[4]))
}

func calcMaxYPos(minX int, maxX int, minY int, maxY int) int {
	canStopOnTarget := false
	// startXVelocity := 0
	distance := 0
	for numSteps := 1; ; numSteps++ {
		distance += numSteps
		if distance >= minX && distance <= maxX {
			canStopOnTarget = true
			// startXVelocity = numSteps
			break
		} else if distance > maxX {
			break
		}
	}
	if !canStopOnTarget {
		panic("Unknown case... can't stop...")
	}

	startYVelocity := -minY - 1
	// fmt.Println(startXVelocity, startYVelocity)

	maxYPos := startYVelocity * (startYVelocity + 1) / 2
	return maxYPos
}

type velocity struct {
	x int
	y int
}

type velocitySlice []velocity

func (vs velocitySlice) Len() int {
	return len(vs)
}

func (vs velocitySlice) Less(i, j int) bool {
	if vs[i].x < vs[j].x {
		return true
	} else if vs[i].x > vs[j].x {
		return false
	} else {
		return vs[i].y < vs[j].y
	}
}

func (vs velocitySlice) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	// isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	var minX, maxX, minY, maxY int
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &minX, &maxX, &minY, &maxY)
	})

	if maxY >= 0 {
		panic("Unknown case... target not below...")
	}

	maxYPos := calcMaxYPos(minX, maxX, minY, maxY)
	fmt.Println("Max y pos:", maxYPos)

	yStepCountVelocities := map[int][]int{}
	for startYVel := minY; startYVel <= 0; startYVel++ {
		positiveOffset := 0
		if startYVel < 0 {
			positiveOffset = (-startYVel-1)*2 + 1
		}

		for numSteps, y, yVel := 0, 0, startYVel; y >= minY; numSteps, y, yVel = numSteps+1, y+yVel, yVel-1 {
			if y <= maxY {
				yStepCountVelocities[numSteps] = append(yStepCountVelocities[numSteps], startYVel)
				if positiveOffset > 0 {
					yStepCountVelocities[numSteps+positiveOffset] = append(yStepCountVelocities[numSteps+positiveOffset], (-startYVel - 1))
				}
			}
		}
	}

	xStepCountVelocities := map[int][]int{}
	stopValues := []int{}
	stopValueVelocities := []int{}
	for startXVel := 0; startXVel <= maxX; startXVel++ {
		for numSteps, x, xVel := 0, 0, startXVel; x <= maxX; numSteps, x, xVel = numSteps+1, x+xVel, xVel-1 {
			if xVel == 0 {
				if x >= minX {
					stopValues = append(stopValues, numSteps)
					stopValueVelocities = append(stopValueVelocities, startXVel)
				}
				break
			}
			if x >= minX {
				// fmt.Println(startXVel, numSteps)
				xStepCountVelocities[numSteps] = append(xStepCountVelocities[numSteps], startXVel)
			}
		}
	}

	velocities := map[velocity]bool{}
	for value, yVels := range yStepCountVelocities {
		xVels := xStepCountVelocities[value]
		for i, stopValue := range stopValues {
			if value >= stopValue {
				xVels = append(xVels, stopValueVelocities[i])
			}
		}
		for _, yv := range yVels {
			for _, xv := range xVels {
				velocities[velocity{xv, yv}] = true
			}
		}
	}

	totalCount := len(velocities)
	// vels := velocitySlice{}
	// for vel := range velocities {
	// 	vels = append(vels, vel)
	// }
	// sort.Sort(vels)
	// for _, v := range vels {
	// 	fmt.Printf("%d,%d\n", v.x, v.y)
	// }

	// fmt.Println(velocities)
	fmt.Println(totalCount)
}
