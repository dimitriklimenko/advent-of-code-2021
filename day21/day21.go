package main

import (
	"bufio"
	"flag"
	"fmt"
	"regexp"
	"strconv"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

var lineRegex = regexp.MustCompile("^Player ([12]) starting position: ([0-9]+)$")

var rollCounts [10]int

func makeRollCounts() {
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				rollCounts[i+j+k]++
			}
		}
	}
}

func scanInput(scanner *bufio.Scanner, p1 *int, p2 *int) {
	for i := 1; i <= 2; i++ {
		scanner.Scan()
		text := scanner.Text()
		bytes := []byte(text)
		matches := lineRegex.FindSubmatch(bytes)
		if len(matches) != 3 {
			panic("Invalid input!?")
		}
		playerNumber, _ := strconv.Atoi(string(matches[1]))
		if playerNumber == 1 {
			*p1, _ = strconv.Atoi(string(matches[2]))
		} else {
			*p2, _ = strconv.Atoi(string(matches[2]))
		}
	}
}

func calcScoreDelta(posn int) int {
	delta := posn
	if delta == 0 {
		delta = 10
	}
	return delta
}

func simulatePart1(positions [2]int) ([2]int, int) {
	var scores [2]int
	delta := 6
	currentPlayer := 0
	numTurns := 0
	for scores[0] < 1000 && scores[1] < 1000 {
		posn := positions[currentPlayer]
		posn = (posn + delta) % 10
		positions[currentPlayer] = posn
		scores[currentPlayer] += calcScoreDelta(posn)

		delta = (delta + 9) % 10
		currentPlayer = 1 - currentPlayer
		numTurns++
	}
	return scores, numTurns
}

type stateCount [21][21][2][10][10]int

func calcPart2(startPos [2]int) [2]int {
	counts := stateCount{}
	counts[0][0][0][startPos[0]][startPos[1]] = 1

	winCounts := [2]int{}

	for score0 := 0; score0 < 21; score0++ {
		for score1 := 0; score1 < 21; score1++ {
			for currentPlayer := 0; currentPlayer < 2; currentPlayer++ {
				for pos0 := 0; pos0 < 10; pos0++ {
					for pos1 := 0; pos1 < 10; pos1++ {
						stateCount := counts[score0][score1][currentPlayer][pos0][pos1]
						if stateCount == 0 {
							continue
						}
						for totalRoll, rollCount := range rollCounts {
							if rollCount == 0 {
								continue
							}
							newScore0 := score0
							newScore1 := score1
							newPos0 := pos0
							newPos1 := pos1
							delta := stateCount * rollCount
							if currentPlayer == 0 {
								newPos0 = (newPos0 + totalRoll) % 10
								newScore0 += calcScoreDelta(newPos0)
								if newScore0 >= 21 {
									winCounts[0] += delta
									continue
								}
							} else {
								newPos1 = (newPos1 + totalRoll) % 10
								newScore1 += calcScoreDelta(newPos1)
								if newScore1 >= 21 {
									winCounts[1] += delta
									continue
								}
							}
							counts[newScore0][newScore1][1-currentPlayer][newPos0][newPos1] += delta
						}
					}
				}
			}
		}
	}
	return winCounts
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	var pos [2]int
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &pos[0], &pos[1])
	})

	pos[0] %= 10
	pos[1] %= 10

	if !*isPart2 {
		scores, numTurns := simulatePart1(pos)
		losingScore := scores[0]
		if scores[1] < scores[0] {
			losingScore = scores[1]
		}
		fmt.Println(losingScore * numTurns * 3)
	} else {
		makeRollCounts()
		universeCounts := calcPart2(pos)
		maxCount := universeCounts[0]
		if universeCounts[1] > universeCounts[0] {
			maxCount = universeCounts[1]
		}
		fmt.Println(maxCount)
	}
}
