package main

import (
	"flag"
	"fmt"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

func countIncreases(depths []int, windowWidth int) int {
	numIncreases := 0
	for i := windowWidth; i < len(depths); i++ {
		if depths[i] > depths[i-windowWidth] {
			numIncreases += 1
		}
	}
	return numIncreases
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	windowWidth := flag.Int("w", 3, "width of sliding window")
	flag.Parse()

	depths := common.ParseInts(filePath)
	numIncreases := countIncreases(depths, *windowWidth)
	fmt.Println(numIncreases)
}
