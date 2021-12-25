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

var digits = map[string]rune{
	"abcefg":  '0',
	"cf":      '1',
	"acdeg":   '2',
	"acdfg":   '3',
	"bcdf":    '4',
	"abdfg":   '5',
	"abdefg":  '6',
	"acf":     '7',
	"abcdefg": '8',
	"abcdfg":  '9',
}

type patterns [10]string
type outputs [4]string

func Contains(s string, c rune) bool {
	for _, c2 := range s {
		if c == c2 {
			return true
		}
	}
	return false
}

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func processLine(text string) string {
	var pats patterns
	var outs outputs
	tokens := strings.Split(text, "|")
	for i, token := range strings.Fields(tokens[0]) {
		pats[i] = token
	}
	for i, token := range strings.Fields(tokens[1]) {
		outs[i] = token
	}
	return solve(pats, outs)
}

func translate(pat string, segMap map[rune]rune) string {
	result := ""
	for _, c := range pat {
		result += string(segMap[c])
	}
	return SortString(result)
}

func solve(pats patterns, outs outputs) string {
	counts := map[rune]int{}
	sizeMap := map[int][]string{}
	for _, pat := range pats {
		sz := len(pat)
		sizeMap[sz] = append(sizeMap[sz], pat)
		for _, c := range pat {
			counts[c]++
		}
	}

	segsFour := sizeMap[4][0]

	segMap := map[rune]rune{}

	for seg, count := range counts {
		if count == 4 {
			segMap[seg] = 'e'
		} else if count == 6 {
			segMap[seg] = 'b'
		} else if count == 7 {
			if Contains(segsFour, seg) {
				segMap[seg] = 'd'
			} else {
				segMap[seg] = 'g'
			}
		} else if count == 8 {
			if Contains(segsFour, seg) {
				segMap[seg] = 'c'
			} else {
				segMap[seg] = 'a'
			}
		} else if count == 9 {
			segMap[seg] = 'f'
		} else {
			panic(fmt.Sprintf("unexpected count value: %d", count))
		}
	}

	result := ""
	for _, out := range outs {
		result += string(digits[translate(out, segMap)])
	}
	return result
}

func processInput(scanner *bufio.Scanner, result *int, addResults bool) {
	*result = 0
	for scanner.Scan() {
		stringValue := processLine(scanner.Text())
		if addResults {
			value, _ := strconv.Atoi(stringValue)
			*result += value
		} else {
			for _, c := range stringValue {
				if c == '1' || c == '4' || c == '7' || c == '8' {
					*result++
				}
			}
		}
	}
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	addResults := flag.Bool("a", true, "add results together")
	flag.Parse()

	var count int
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		processInput(scanner, &count, *addResults)
	})
	fmt.Println(count)
}
