package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

type edge struct {
	start string
	end   string
}

type edgeMap map[string][]string
type visitMap map[string]bool

func scanInput(scanner *bufio.Scanner, edges *[]edge) {
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), "-")
		*edges = append(*edges, edge{tokens[0], tokens[1]})
	}
}

func countPaths(em edgeMap, current string, visited visitMap, canVisitTwice bool) int {
	paths := 0
	for _, nb := range em[current] {
		newVisited := visitMap{}
		for k, v := range visited {
			newVisited[k] = v
		}
		if current[0] >= 'a' {
			newVisited[current] = true
		}

		if nb == "end" {
			paths += 1
		} else if !visited[nb] {
			paths += countPaths(em, nb, newVisited, canVisitTwice)
		} else {
			if nb != "start" && canVisitTwice {
				paths += countPaths(em, nb, newVisited, false)
			}
		}
	}
	return paths
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	canVisitTwice := flag.Bool("", true, "can visit a small cave twice")
	flag.Parse()

	edges := []edge{}
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &edges)
	})
	em := edgeMap{}
	for _, e := range edges {
		em[e.start] = append(em[e.start], e.end)
		em[e.end] = append(em[e.end], e.start)
	}

	numPaths := countPaths(em, "start", visitMap{}, *canVisitTwice)
	fmt.Println(numPaths)
}
