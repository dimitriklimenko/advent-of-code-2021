package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type result struct {
	left    int
	right   int
	newChar int
}

func scanInput(scanner *bufio.Scanner, template *string, sources *[]string, insertions *[]string) {
	scanner.Scan()
	*template = scanner.Text()
	scanner.Scan()

	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), "->")
		*sources = append(*sources, strings.TrimSpace(tokens[0]))
		*insertions = append(*insertions, strings.TrimSpace(tokens[1]))
	}
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	numSteps := flag.Int("", 40, "number of steps to iterate")
	flag.Parse()

	var polymer string
	var sources []string
	var insertions []string
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &polymer, &sources, &insertions)
	})

	allChars := map[byte]bool{}
	for i, s := range sources {
		allChars[s[0]] = true
		allChars[s[1]] = true
		allChars[insertions[i][0]] = true
	}

	charIndices := map[byte]int{}
	i := 0
	for c, _ := range allChars {
		charIndices[c] = i
		i++
	}

	pairIndices := map[string]int{}
	for i, s := range sources {
		pairIndices[s] = i
	}

	var results []result
	for i, s := range sources {
		resultString := string(s[0]) + insertions[i] + string(s[1])
		results = append(results, result{
			left:    pairIndices[resultString[0:2]],
			right:   pairIndices[resultString[1:3]],
			newChar: charIndices[insertions[i][0]],
		})
	}

	pairCounts := make([]int, len(sources))
	charCounts := make([]int, len(allChars))

	for i := 0; i < len(polymer)-1; i++ {
		pairCounts[pairIndices[polymer[i:i+2]]]++
	}
	for i := range polymer {
		charCounts[charIndices[polymer[i]]]++
	}

	for i := 0; i < *numSteps; i++ {
		newPairCounts := make([]int, len(sources))
		for j, pc := range pairCounts {
			r := results[j]
			newPairCounts[r.left] += pc
			newPairCounts[r.right] += pc
			charCounts[r.newChar] += pc
		}
		pairCounts = newPairCounts
	}

	minCount := charCounts[0]
	maxCount := charCounts[0]
	for _, count := range charCounts[1:] {
		if count < minCount {
			minCount = count
		}
		if count > maxCount {
			maxCount = count
		}
	}
	fmt.Println(maxCount - minCount)
}
