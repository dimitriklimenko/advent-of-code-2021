package main

import (
	"bufio"
	"flag"
	"fmt"
	"sort"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

func scanInput(scanner *bufio.Scanner, lines *[]string) {
	for scanner.Scan() {
		*lines = append(*lines, scanner.Text())
	}
}

func getCloser(opener rune) rune {
	switch opener {
	case '(':
		return ')'
	case '[':
		return ']'
	case '{':
		return '}'
	case '<':
		return '>'
	default:
		return '0'
	}
}

func checkLine(line string) (rune, string) {
	stack := []rune{}
	for _, c := range line {
		if c == '(' || c == '[' || c == '{' || c == '<' {
			stack = append(stack, c)
		} else {
			if len(stack) <= 0 {
				return c, ""
			}
			closer := getCloser(stack[len(stack)-1])
			if c != closer {
				return c, ""
			}
			stack = stack[:len(stack)-1]
		}
	}
	completionString := ""
	for i := len(stack) - 1; i >= 0; i-- {
		completionString += string(getCloser(stack[i]))
	}
	return '0', completionString
}

func getScore(completionString string) int {
	total := 0
	for _, c := range completionString {
		total *= 5
		switch c {
		case ')':
			total += 1
		case ']':
			total += 2
		case '}':
			total += 3
		case '>':
			total += 4
		}
	}
	return total
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	flag.Parse()

	lines := []string{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &lines)
	})

	score := 0
	completionScores := []int{}
	for _, line := range lines {
		illegalChar, completionString := checkLine(line)
		if completionString != "" {
			completionScores = append(completionScores, getScore(completionString))
		}
		switch illegalChar {
		case ')':
			score += 3
		case ']':
			score += 57
		case '}':
			score += 1197
		case '>':
			score += 25137
		}
	}
	sort.Ints(completionScores)
	median := completionScores[len(completionScores)/2]

	fmt.Println(score)
	fmt.Println(median)
}
