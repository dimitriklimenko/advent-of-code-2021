package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type BingoBoard [5][5]int
type TimeBoard [5][5]int

func readBoard(scanner *bufio.Scanner, board *BingoBoard) bool {
	if !scanner.Scan() {
		return false
	}
	for i := 0; i < 5; i++ {
		scanner.Scan()
		for j, token := range strings.Fields(scanner.Text()) {
			board[i][j], _ = strconv.Atoi(token)
		}
	}
	return true
}

func scanInput(scanner *bufio.Scanner, drawOrder *[]int, drawOrderMap map[int]int, boards *[]BingoBoard) {
	scanner.Scan()

	for index, token := range strings.Split(scanner.Text(), ",") {
		number, _ := strconv.Atoi(token)
		*drawOrder = append(*drawOrder, number)
		drawOrderMap[number] = index
	}

	var board BingoBoard
	for readBoard(scanner, &board) {
		*boards = append(*boards, board)
	}
}

func makeTimeBoard(drawOrderMap map[int]int, board BingoBoard) TimeBoard {
	timeBoard := TimeBoard{}
	for i, row := range board {
		for j, number := range row {
			var ok bool
			timeBoard[i][j], ok = drawOrderMap[number]
			if !ok {
				panic("Number on board not present in draw order!?")
			}
		}
	}
	return timeBoard
}

func calcWinTime(timeBoard TimeBoard) int {
	bestWinTime := -1
	for _, isRows := range []bool{true, false} {
		for i := 0; i < 5; i++ {
			winTime := 0
			for j := 0; j < 5; j++ {
				var time int
				if isRows {
					time = timeBoard[i][j]
				} else {
					time = timeBoard[j][i]
				}

				if time > winTime {
					winTime = time
				}
			}
			if bestWinTime == -1 || winTime < bestWinTime {
				bestWinTime = winTime
			}
		}
	}
	return bestWinTime
}

func calcScore(board BingoBoard, timeBoard TimeBoard, winTime int, winNumber int) int {
	unmarkedSum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			number := board[i][j]
			time := timeBoard[i][j]
			if time > winTime {
				unmarkedSum += number
			}
		}
	}

	return unmarkedSum * winNumber
}

func calcFinalScore(drawOrder []int, drawOrderMap map[int]int, boards []BingoBoard, useLastBoard bool) int {
	timeBoards := []TimeBoard{}
	bestWinTime := -1
	bestBoardIdx := -1

	for i, board := range boards {
		timeBoard := makeTimeBoard(drawOrderMap, board)
		timeBoards = append(timeBoards, timeBoard)

		winTime := calcWinTime(timeBoard)
		if bestWinTime == -1 ||
			(!useLastBoard && winTime < bestWinTime) ||
			(useLastBoard && winTime > bestWinTime) {
			bestWinTime = winTime
			bestBoardIdx = i
		}
	}

	return calcScore(boards[bestBoardIdx], timeBoards[bestBoardIdx], bestWinTime, drawOrder[bestWinTime])
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	useLastBoard := flag.Bool("last", true, "set to use last board instead of first board")
	flag.Parse()

	boards := []BingoBoard{}
	drawOrder := []int{}
	drawOrderMap := map[int]int{}

	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &drawOrder, drawOrderMap, &boards)
	})
	score := calcFinalScore(drawOrder, drawOrderMap, boards, *useLastBoard)
	fmt.Println(score)
}
