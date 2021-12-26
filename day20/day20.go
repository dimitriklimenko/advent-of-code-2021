package main

import (
	"bufio"
	"flag"
	"fmt"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

const (
	dark      uint8 = 0
	darkChar        = '.'
	light     uint8 = 1
	lightChar       = '#'
)

type algorithm []uint8
type image [][]uint8

func printImage(im image) {
	for _, row := range im {
		for _, col := range row {
			if col == dark {
				fmt.Printf("%c", darkChar)
			} else {
				fmt.Printf("%c", lightChar)
			}
		}
		fmt.Println()
	}
}

func parseRow(text string) []uint8 {
	row := []uint8{}
	for _, c := range text {
		if c == darkChar {
			row = append(row, dark)
		} else if c == lightChar {
			row = append(row, light)
		} else {
			panic(fmt.Sprintf("Unexpected pixel: %c", c))
		}
	}
	return row
}

func scanInput(scanner *bufio.Scanner, algo *algorithm, im *image) {
	scanner.Scan()
	*algo = parseRow(scanner.Text())
	scanner.Scan()
	for scanner.Scan() {
		*im = append(*im, parseRow(scanner.Text()))
	}
}

func makeBackgroundImage(background uint8, numRows int, numCols int) image {
	newImage := make(image, numRows)
	for i := range newImage {
		newImage[i] = make([]uint8, numCols)
		if background != 0 {
			for j := 0; j < numCols; j++ {
				newImage[i][j] = background
			}
		}
	}
	return newImage
}

func makeBiggerCopy(im image, background uint8) image {
	newImage := makeBackgroundImage(background, len(im)+2, len(im[0])+2)
	for i, row := range im {
		for j, px := range row {
			newImage[i+1][j+1] = px
		}
	}
	return newImage
}

func enhance(im image, background uint8, algo algorithm) (image, uint8) {
	if background == dark {
		background = algo[0]
	} else {
		background = algo[511]
	}
	numRows := len(im)
	numCols := len(im[0])
	newImage := makeBackgroundImage(background, numRows+2, numCols+2)
	for i := 0; i < numRows-2; i++ {
		for j := 0; j < numCols-2; j++ {
			index := 0
			for di := 0; di < 3; di++ {
				for dj := 0; dj < 3; dj++ {
					index = index*2 + int(im[i+di][j+dj])
				}
			}
			result := algo[index]
			newImage[i+2][j+2] = result
		}
	}
	return newImage, background
}

func countLitPixes(im image, background uint8) int {
	if background == light {
		panic("Infinite pixels!")
	}
	count := 0
	for _, row := range im {
		for _, px := range row {
			if px == 1 {
				count++
			}
		}
	}
	return count
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	numEnhancements := flag.Int("n", 50, "# of times to enhance")
	flag.Parse()
	var algo algorithm
	var im image
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &algo, &im)
	})
	if len(algo) != 512 {
		panic("Wrong algorithm length!")
	}

	// tart
	var background uint8 = 0
	im = makeBiggerCopy(makeBiggerCopy(im, background), background)
	// printImage(im)
	// fmt.Println()
	for i := 0; i < *numEnhancements; i++ {
		im, background = enhance(im, background, algo)
		// printImage(im)
		// fmt.Println()
	}
	count := countLitPixes(im, background)
	fmt.Println(count)
}
