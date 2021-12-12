package main

import (
	"flag"
	"fmt"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

func calcPowerConsumption(report []string) int {
	numEntries := len(report)
	// Assume all numbers have the same # of bits.
	numBits := len(report[0])

	zeroCounts := make([]int, numBits)
	for _, number := range report {
		for i, bit := range number {
			if bit == '0' {
				zeroCounts[i]++
			}
		}
	}

	placeValue := 1
	gammaRate := 0
	epsilonRate := 0
	for i := numBits - 1; i >= 0; i-- {
		if zeroCounts[i]*2 == numEntries {
			panic("Tie for most common bit; result is ill-defined!")
		} else if zeroCounts[i]*2 < numEntries {
			epsilonRate += placeValue
		} else {
			gammaRate += placeValue
		}
		placeValue *= 2
	}
	return gammaRate * epsilonRate
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	flag.Parse()

	report := common.ParseStrings(filePath)
	powerConsumption := calcPowerConsumption(report)
	fmt.Println(powerConsumption)
}
