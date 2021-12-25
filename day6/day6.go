package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type fishCounts [9]int

func scanInput(scanner *bufio.Scanner, counts *fishCounts) {
	scanner.Scan()
	for _, token := range strings.Split(scanner.Text(), ",") {
		timer, _ := strconv.Atoi(token)
		counts[timer] += 1
	}
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	numDays := flag.Int("d", 256, "number of days to simulate")
	flag.Parse()

	counts := fishCounts{}

	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &counts)
	})

	for time := 0; time < *numDays; time++ {
		numNew := counts[0]

		for i := 0; i <= 7; i++ {
			counts[i] = counts[i+1]
		}

		counts[8] = numNew
		counts[6] += numNew
	}

	total := 0
	for _, count := range counts {
		total += count
	}

	fmt.Println(total)
}
