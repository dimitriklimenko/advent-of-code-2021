package main

import (
	"flag"
	"fmt"
	"strconv"

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

func filterByBit(report []string, index int, keepMostCommon bool) []string {
	zeros := []string{}
	ones := []string{}
	for _, number := range report {
		if number[index] == '0' {
			zeros = append(zeros, number)
		} else {
			ones = append(ones, number)
		}
	}

	if len(zeros) > len(ones) {
		if keepMostCommon {
			return zeros
		} else {
			return ones
		}
	} else {
		if keepMostCommon {
			return ones
		} else {
			return zeros
		}
	}
}

func calcRating(report []string, keepMostCommon bool) uint64 {
	// Assume all numbers have the same # of bits.
	numBits := len(report[0])

	currentIndex := 0

	for len(report) > 1 {
		if currentIndex >= numBits {
			panic("Non-unique numbers; could not determine a rating!")
		}
		report = filterByBit(report, currentIndex, keepMostCommon)
		currentIndex += 1
	}

	if len(report) == 0 {
		panic("Least common value occurred zero times!")
	}

	rating, _ := strconv.ParseUint(report[0], 2, 64)
	return rating
}

func calcLifeSupportRating(report []string) uint64 {
	return calcRating(report, true) * calcRating(report, false)
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	isLifeSupport := flag.Bool("life", true, "set to calculat life support rating")
	flag.Parse()

	report := common.ParseStrings(filePath)
	if *isLifeSupport {
		lifeSupportRating := calcLifeSupportRating(report)
		fmt.Println(lifeSupportRating)
	} else {
		powerConsumption := calcPowerConsumption(report)
		fmt.Println(powerConsumption)
	}
}
