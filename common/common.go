package common

import (
	"bufio"
	"os"
	"strconv"
)

func ParseFile(filePath *string, scanFunc func(*bufio.Scanner)) {
	// Default to stdin, so we can pipe in input more conveniently if needed.
	var fileHandle *os.File = os.Stdin
	var err error
	if *filePath != "" {
		fileHandle, err = os.Open(*filePath)
		if err != nil {
			panic(err)
		}
		defer fileHandle.Close()
	}
	scanner := bufio.NewScanner(fileHandle)
	scanFunc(scanner)
}

func ScanInts(scanner *bufio.Scanner, output *[]int) {
	for scanner.Scan() {
		num, _ := strconv.Atoi(scanner.Text())
		*output = append(*output, num)
	}
}

func ParseInts(filePath *string) []int {
	numbers := []int{}
	ParseFile(filePath, func(scanner *bufio.Scanner) {
		ScanInts(scanner, &numbers)
	})
	return numbers
}

func ScanStrings(scanner *bufio.Scanner, output *[]string) {
	for scanner.Scan() {
		token := scanner.Text()
		*output = append(*output, token)
	}
}

func ParseStrings(filePath *string) []string {
	strings := []string{}
	ParseFile(filePath, func(scanner *bufio.Scanner) {
		ScanStrings(scanner, &strings)
	})
	return strings
}
