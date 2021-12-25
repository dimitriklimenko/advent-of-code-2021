package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type snumber struct {
	number int
	left   *snumber
	right  *snumber
}

func printSnumber(sn *snumber) {
	if sn.left == nil {
		fmt.Print(sn.number)
	} else {
		fmt.Print("[")
		printSnumber(sn.left)
		fmt.Print(",")
		printSnumber(sn.right)
		fmt.Print("]")
	}
}

func scanInput(scanner *bufio.Scanner, lines *[]string) {
	for scanner.Scan() {
		*lines = append(*lines, scanner.Text())
	}
}

func readSnumber(s string, index int) (snumber, int) {
	sn := snumber{}
	startIndex := index
	if s[index] == '[' {
		index += 1
		left, offset := readSnumber(s, index)
		index += offset
		if s[index] != ',' {
			panic("Invalid snumber!")
		}
		index += 1
		right, offset2 := readSnumber(s, index)
		index += offset2
		if s[index] != ']' {
			panic("Invalid snumber!")
		}
		index += 1
		sn.left = &left
		sn.right = &right
	} else {
		for ; ; index++ {
			if s[index] == ',' || s[index] == ']' {
				break
			}
		}
		value, _ := strconv.Atoi(s[startIndex:index])
		sn.number = value
	}
	return sn, index - startIndex
}

func add(sn1 *snumber, sn2 *snumber) snumber {
	result := snumber{
		left:  sn1,
		right: sn2,
	}
	// printSnumber(&result)
	// fmt.Println()
	reduce(&result)
	return result
}

func copy(sn *snumber) snumber {
	if sn.left == nil {
		return snumber{number: sn.number}
	} else {
		left := copy(sn.left)
		right := copy(sn.right)
		return snumber{left: &left, right: &right}
	}
}

func reduce(sn *snumber) {
	for {
		exploded, _, _ := explode(sn, 0)
		if exploded {
			// printSnumber(sn)
			// fmt.Println()
			continue
		}
		if split(sn) {
			// printSnumber(sn)
			// fmt.Println()
			continue
		}
		return
	}
}

func addLeft(sn *snumber, delta int) int {
	if sn.left == nil {
		sn.number += delta
		return 0
	} else {
		result := addLeft(sn.left, delta)
		if result == 0 {
			return 0
		}
		result = addRight(sn.right, delta)
		return result
	}
}

func addRight(sn *snumber, delta int) int {
	if sn.left == nil {
		sn.number += delta
		return 0
	} else {
		result := addRight(sn.right, delta)
		if result == 0 {
			return 0
		}
		result = addRight(sn.left, delta)
		return result
	}
}

func explode(sn *snumber, depth int) (bool, int, int) {
	if sn.left == nil {
		// Regular numbers can't be exploded
		return false, 0, 0
	}
	if sn.left != nil {
		// If both numbers in the pair are regular numbers and the depth is sufficient, we explode.
		if depth >= 4 && sn.left.left == nil && sn.right.right == nil {
			leftVal := sn.left.number
			rightVal := sn.right.number
			sn.left = nil
			sn.right = nil
			sn.number = 0
			return true, leftVal, rightVal
		}

		// Check for explosion in the left child.
		exploded, leftAdd, rightAdd := explode(sn.left, depth+1)
		if exploded {
			if rightAdd > 0 {
				rightAdd = addLeft(sn.right, rightAdd)
			}
			return true, leftAdd, rightAdd
		}

		// Check for explosion in the right child.
		exploded, leftAdd, rightAdd = explode(sn.right, depth+1)
		if exploded {
			if leftAdd > 0 {
				leftAdd = addRight(sn.left, leftAdd)
			}
		}
		return exploded, leftAdd, rightAdd
	}
	return false, 0, 0
}

func split(sn *snumber) bool {
	if sn.left == nil {
		if sn.number >= 10 {
			splitValue := sn.number / 2
			sn.left = &snumber{number: splitValue}
			sn.right = &snumber{number: splitValue + sn.number%2}
			sn.number = 0
			return true
		}
		return false
	} else {
		if split(sn.left) {
			return true
		}
		return split(sn.right)
	}
}

func magnitude(sn *snumber) int {
	if sn.left == nil {
		return sn.number
	} else {
		return 3*magnitude(sn.left) + 2*magnitude(sn.right)
	}
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	// isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	lines := []string{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &lines)
	})

	sns := []*snumber{}
	for _, line := range lines {
		sn, _ := readSnumber(line, 0)
		sns = append(sns, &sn)
	}

	maxMagnitude := 0
	for i := range sns {
		for j := range sns {
			if i == j {
				continue
			}
			a := copy(sns[i])
			b := copy(sns[j])
			sum := add(&a, &b)
			mag := magnitude(&sum)
			if mag > maxMagnitude {
				maxMagnitude = mag
			}
		}
	}

	result := sns[0]
	for _, sn := range sns[1:] {
		newResult := add(result, sn)
		result = &newResult
	}
	mag := magnitude(result)
	fmt.Println(mag)

	fmt.Println(maxMagnitude)

	// printSnumber(result)
	// fmt.Println()
}
