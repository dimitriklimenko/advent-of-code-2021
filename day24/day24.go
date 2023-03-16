package main

import (
	"fmt"
	"strconv"
)

var number = "12345678912345"

func equal(a int, b int) int {
	if a == b {
		return 1
	} else {
		return 0
	}
}

func notEqual(a int, b int) int {
	if a != b {
		return 1
	} else {
		return 0
	}
}

func main() {
	ws := []int{}
	for _, c := range number {
		val, _ := strconv.Atoi(string(c))
		ws = append(ws, val)
	}
	fmt.Println()

	var w, x, y, z int

	w = ws[0]
	x = (w+8)%26 + 15

	w = ws[1]
	x = notEqual(x, w)
	y = (w + 11) * x
	z = y

	w = ws[2]
	x = notEqual((z%26)+13, w)
	y = (w + 2) * x
	z = y

}
