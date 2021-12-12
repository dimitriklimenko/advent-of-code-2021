package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func readInts(scanner *bufio.Scanner) []int {
	numbers := []int{}
	for scanner.Scan() {
		num, _ := strconv.Atoi(scanner.Text())
		numbers = append(numbers, num)
	}
	return numbers
}

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

	var fileHandle *os.File = os.Stdin
	var err error
	if *filePath != "" {
		fileHandle, err = os.Open(*filePath)
		if err != nil {
			panic(err)
		}
		defer fileHandle.Close()
	}
	depths := readInts(bufio.NewScanner(fileHandle))
	numIncreases := countIncreases(depths, *windowWidth)
	fmt.Println(numIncreases)
}
