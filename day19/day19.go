package main

import (
	"bufio"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type vector struct {
	x int
	y int
	z int
}

type report []vector

type lookupItem struct {
	col int
	sgn int
}

var rotationLookup = [3][3]lookupItem{
	{
		{},
		{2, +1},
		{1, -1},
	},
	{
		{2, -1},
		{},
		{0, +1},
	},
	{
		{1, +1},
		{0, -1},
		{},
	},
}

type rotationMatrix [3][3]int

var rotationMatrices []rotationMatrix

func makeRotationMatrices() {
	for topCol := 0; topCol < 3; topCol++ {
		for _, topVal := range []int{+1, -1} {
			for middleCol := 0; middleCol < 3; middleCol++ {
				if middleCol == topCol {
					continue
				}
				for _, middleVal := range []int{+1, -1} {
					matrix := rotationMatrix{}
					matrix[0][topCol] = topVal
					matrix[1][middleCol] = middleVal
					lookupResult := rotationLookup[topCol][middleCol]
					matrix[2][lookupResult.col] = topVal * middleVal * lookupResult.sgn
					rotationMatrices = append(rotationMatrices, matrix)
				}
			}
		}
	}
}

func rotateVec(vec vector, rot rotationMatrix) vector {
	return vector{
		rot[0][0]*vec.x + rot[0][1]*vec.y + rot[0][2]*vec.z,
		rot[1][0]*vec.x + rot[1][1]*vec.y + rot[1][2]*vec.z,
		rot[2][0]*vec.x + rot[2][1]*vec.y + rot[2][2]*vec.z,
	}
}

func copyReport(rep report) report {
	result := report{}
	result = append(result, rep...)
	return result
}

func rotateReport(rep *report, rot *rotationMatrix) {
	for i := range *rep {
		(*rep)[i] = rotateVec((*rep)[i], *rot)
	}
}

func translateReport(rep *report, offset vector) {
	for i := range *rep {
		(*rep)[i].x += offset.x
		(*rep)[i].y += offset.y
		(*rep)[i].z += offset.z
	}
}

var headerRegex = regexp.MustCompile("^--- scanner ([0-9]+) ---$")

func readVector(text string) vector {
	tokens := strings.Split(text, ",")
	x, _ := strconv.Atoi(tokens[0])
	y, _ := strconv.Atoi(tokens[1])
	z, _ := strconv.Atoi(tokens[2])
	return vector{x, y, z}
}

func scanInput(scanner *bufio.Scanner, reports *[]report) {
	for reportNumber := 0; scanner.Scan(); reportNumber++ {
		headerText := scanner.Text()
		matches := headerRegex.FindSubmatch([]byte(headerText))
		if len(matches) != 2 {
			panic(fmt.Sprintf("Invalid header: %s", headerText))
		}
		numberRead, _ := strconv.Atoi(string(matches[1]))
		if numberRead != reportNumber {
			panic(fmt.Sprintf("Wrong report number; expected %d but got %d", reportNumber, numberRead))
		}

		rep := report{}
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				break
			}

			rep = append(rep, readVector(text))
		}
		*reports = append(*reports, rep)
	}
}

func countOverlap(r1 report, r2 report) int {
	r1Map := map[vector]bool{}
	for _, vec := range r1 {
		r1Map[vec] = true
	}

	count := 0
	for _, vec := range r2 {
		if r1Map[vec] {
			count++
		}
	}
	// if count >= 12 {
	// 	for _, vec := range r2 {
	// 		if r1Map[vec] {
	// 			fmt.Printf("%d,%d,%d\n", vec.x, vec.y, vec.z)
	// 		}
	// 	}
	// 	fmt.Println()
	// }
	return count
}

func subtract(v1 vector, v2 vector) vector {
	return vector{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
	}
}

func abs(x int) int {
	if x >= 0 {
		return x
	} else {
		return -x
	}
}

func l1Norm(v1 vector) int {
	return abs(v1.x) + abs(v1.y) + abs(v1.z)
}

func findOverlaps(r1 report, r2 report) (*rotationMatrix, vector) {
	for _, rot := range rotationMatrices {
		rotatedR2 := copyReport(r2)
		rotateReport(&rotatedR2, &rot)
		for i := 0; i < len(r1); i++ {
			for j := 0; j < len(r2); j++ {
				offset := subtract(r1[i], rotatedR2[j])
				translateReport(&rotatedR2, offset)
				count := countOverlap(r1, rotatedR2)
				if count >= 12 {
					trueOffset := subtract(r1[i], rotateVec(r2[j], rot))
					return &rot, trueOffset
				}
			}
		}
	}
	return nil, vector{}
}

func transformRegions(reports []report) ([]*rotationMatrix, []vector) {
	numReports := len(reports)
	known := []int{0}
	unknown := map[int]bool{}
	for i := 1; i < len(reports); i++ {
		unknown[i] = true
	}

	matrices := make([]*rotationMatrix, len(reports))
	offsets := make([]vector, len(reports))

	fmt.Printf("%02d/%d: %d\n", len(known), numReports, 0)
	for len(known) < len(reports) {
		found := false
		for target := range unknown {
			for _, source := range known {
				rot, offset := findOverlaps(reports[source], reports[target])
				if rot != nil {
					delete(unknown, target)
					known = append(known, target)
					fmt.Printf("%02d/%d: %d\n", len(known), numReports, target)

					rotateReport(&reports[target], rot)
					translateReport(&reports[target], offset)
					matrices[target] = rot
					offsets[target] = offset
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			panic("Could not find an overlap!?")
		}
	}
	fmt.Println()
	return matrices, offsets
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	// isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	makeRotationMatrices()

	reports := []report{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &reports)
	})

	_, offsets := transformRegions(reports)
	allBeacons := map[vector]bool{}
	for _, rep := range reports {
		for _, vec := range rep {
			allBeacons[vec] = true
		}
	}
	fmt.Println("# of beacons:", len(allBeacons))

	maxDistance := 0
	for i := 0; i < len(reports); i++ {
		for j := 0; j < len(reports); j++ {
			if i == j {
				continue
			}
			distance := l1Norm(subtract(offsets[i], offsets[j]))
			if distance > maxDistance {
				maxDistance = distance
			}
		}
	}

	fmt.Println("max distance:", maxDistance)
}
