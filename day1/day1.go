package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	filePath := flag.String("f", "", "file to read from")
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
	scanner := bufio.NewScanner(fileHandle)

	numIncreases := 0

	scanner.Scan()
	prevDepth, _ := strconv.Atoi(scanner.Text())
	for scanner.Scan() {
		depth, _ := strconv.Atoi(scanner.Text())
		if depth > prevDepth {
			numIncreases += 1
		}
		prevDepth = depth
	}
	fmt.Println(numIncreases)
}
