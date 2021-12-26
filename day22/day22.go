package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

var lineRegex = regexp.MustCompile("(on|off) x=(-?[0-9]+)..(-?[0-9]+),y=(-?[0-9]+)..(-?[0-9]+),z=(-?[0-9]+)..(-?[0-9]+)$")

type coordRange struct {
	min int
	max int
}

type cuboid struct {
	ranges [3]coordRange
}

type instruction struct {
	turnOn bool
	target cuboid
}

type limitedGrid [101][101][101]bool

func parseRange(parts [][]byte) coordRange {
	coor := coordRange{}
	coor.min, _ = strconv.Atoi(string(parts[0]))
	coor.max, _ = strconv.Atoi(string(parts[1]))
	if coor.min > coor.max {
		panic("Backwards range...?")
	}
	return coor
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
			target: cuboid{
				ranges: [3]coordRange{
					parseRange(matches[2:4]),
					parseRange(matches[4:6]),
					parseRange(matches[6:8]),
				},
			},
		}
		*steps = append(*steps, step)
	}
}

func max(v1 int, v2 int) int {
	result := v1
	if v2 > result {
		result = v2
	}
	return result
}

func min(v1 int, v2 int) int {
	result := v1
	if v2 < result {
		result = v2
	}
	return result
}

func processStepPart1(grid *limitedGrid, step instruction) {
	for i := max(step.target.ranges[0].min, -50) + 50; i <= min(step.target.ranges[0].max, +50)+50; i++ {
		for j := max(step.target.ranges[1].min, -50) + 50; j <= min(step.target.ranges[1].max, +50)+50; j++ {
			for k := max(step.target.ranges[2].min, -50) + 50; k <= min(step.target.ranges[2].max, +50)+50; k++ {
				(*grid)[i][j][k] = step.turnOn
			}
		}
	}
}

func countCubesOnPart1(grid limitedGrid) int {
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

type cuboidState struct {
	boundingBox *cuboid
	children    []*cuboidState
}

const (
	noOverlap  = 0
	hasOverlap = 1
	contained  = 2
)

func checkRangeOverlap(boundingRange coordRange, target coordRange) int {
	if boundingRange.min <= target.min && boundingRange.max >= target.max {
		return contained
	} else if boundingRange.min <= target.max && boundingRange.max >= target.min {
		return hasOverlap
	} else {
		return noOverlap
	}
}

func calcDifference(boundingRange coordRange, target coordRange) (*coordRange, []coordRange) {
	offRange := coordRange{
		min: max(boundingRange.min, target.min),
		max: min(boundingRange.max, target.max),
	}
	if offRange.min > offRange.max {
		return nil, []coordRange{
			boundingRange,
		}
	}
	onRanges := []coordRange{}
	if boundingRange.min < offRange.min {
		onRanges = append(onRanges, coordRange{boundingRange.min, offRange.min - 1})
	}
	if boundingRange.max > offRange.max {
		onRanges = append(onRanges, coordRange{offRange.max + 1, boundingRange.max})
	}
	return &offRange, onRanges
}

func checkCuboidOverlap(boundingBox cuboid, target cuboid) int {
	result := contained
	for i := 0; i < 3; i++ {
		rangeResult := checkRangeOverlap(boundingBox.ranges[i], target.ranges[i])
		if rangeResult < result {
			result = rangeResult
		}
		if result == noOverlap {
			break
		}
	}
	return result
}

func turnOff(state *cuboidState, target cuboid) bool {
	newChildren := []*cuboidState{}
	if len(state.children) == 0 {
		if state.boundingBox != nil {
			newRanges := state.boundingBox.ranges
			for i := 0; i < 3; i++ {
				offRange, onRanges := calcDifference(state.boundingBox.ranges[i], target.ranges[i])
				if offRange == nil {
					panic("No overlap; shouldn't be here...")
				}
				for _, coor := range onRanges {
					childRanges := newRanges
					childRanges[i] = coor
					newChildren = append(newChildren, &cuboidState{
						boundingBox: &cuboid{childRanges},
					})
				}
				newRanges[i] = *offRange
			}
		}
	} else {
		for _, child := range state.children {
			mustKeep := true
			overlapResult := checkCuboidOverlap(*child.boundingBox, target)
			if overlapResult != noOverlap {
				mustKeep = turnOff(child, target)
			}
			if mustKeep {
				newChildren = append(newChildren, child)
			}
		}
	}
	state.children = newChildren
	return len(newChildren) != 0
}

func turnOn(state *cuboidState, target cuboid) {
	done := false
	newChildren := []*cuboidState{}
	if len(state.children) == 0 {
		// If all the lights in this cuboid are already on, we're done.
		if state.boundingBox != nil {
			return
		}
	} else {
		for _, child := range state.children {
			mustKeep := true
			overlapResult := checkCuboidOverlap(*child.boundingBox, target)
			if overlapResult == contained {
				if !done {
					turnOn(child, target)
					done = true
				} else {
					mustKeep = turnOff(child, target)
				}
			} else if overlapResult == hasOverlap {
				mustKeep = turnOff(child, target)
			}
			if mustKeep {
				newChildren = append(newChildren, child)
			}
		}
	}
	if !done {
		newChildren = append(newChildren, &cuboidState{
			boundingBox: &target,
		})
	}
	state.children = newChildren
}

func processStepPart2(state *cuboidState, step instruction) {
	if step.turnOn {
		turnOn(state, step.target)
	} else {
		turnOff(state, step.target)
	}
}

func countCubesOnPart2(state cuboidState) *big.Int {
	var count big.Int
	if len(state.children) == 0 {
		if state.boundingBox != nil {
			count.SetInt64(1)
			for _, coor := range state.boundingBox.ranges {
				other := big.Int{}
				other.SetUint64(uint64(coor.max - coor.min + 1))
				count.Mul(&count, &other)
			}
		}
	} else {
		for _, child := range state.children {
			count.Add(&count, countCubesOnPart2(*child))
		}
	}
	return &count
}

func main() {
	filePath := flag.String("f", "", "file to read from / write to")
	doGen := flag.Bool("g", false, "generate test case")
	flag.Parse()

	if *doGen {
		generate(filePath)
		// return
	}

	steps := []instruction{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &steps)
	})
	fmt.Println()

	var limitedCount int
	var fullCount *big.Int
	grid := limitedGrid{}
	state := cuboidState{}

	for _, step := range steps {
		// limitedCount = countCubesOnPart1(grid)
		// fullCount = countCubesOnPart2(state)
		// fmt.Printf("%d, %d\n", limitedCount, fullCount)
		// if int64(limitedCount) != fullCount.Int64() {
		// 	print("!")
		// }

		processStepPart1(&grid, step)
		processStepPart2(&state, step)
	}

	limitedCount = countCubesOnPart1(grid)
	fullCount = countCubesOnPart2(state)
	fmt.Printf("%d, %d\n", limitedCount, fullCount)
}

func generateRandomInstruction(rng *rand.Rand) instruction {
	result := instruction{}
	result.turnOn = rng.Int()%2 == 0
	for i := 0; i < 2; i++ {
		v1 := rng.Int()%101 - 50
		v2 := rng.Int()%101 - 50
		result.target.ranges[i].min = min(v1, v2)
		result.target.ranges[i].max = max(v1, v2)
	}
	return result
}

func generate(filePath *string) {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	steps := []instruction{}
	for i := 0; i < 100; i++ {
		steps = append(steps, generateRandomInstruction(rng))
	}

	var fileHandle *os.File = os.Stdout
	var err error
	if *filePath != "" {
		fileHandle, err = os.Create(*filePath)
		if err != nil {
			panic(err)
		}
		defer fileHandle.Close()
	}
	writer := bufio.NewWriter(fileHandle)
	defer writer.Flush()

	for _, step := range steps {
		typeString := "off"
		if step.turnOn {
			typeString = "on"
		}
		stepString := fmt.Sprintf(
			"%s x=%d..%d,y=%d..%d,z=%d..%d\n",
			typeString,
			step.target.ranges[0].min, step.target.ranges[0].max,
			step.target.ranges[1].min, step.target.ranges[1].max,
			step.target.ranges[2].min, step.target.ranges[2].max,
		)
		writer.WriteString(stepString)
	}
}
