package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/dimitriklimenko/advent-of-code-2021/common"
)

const (
	sum         = 0
	product     = 1
	minimum     = 2
	maximum     = 3
	literal     = 4
	greaterThan = 5
	lessThan    = 6
	equalTo     = 7
)

func scanInput(scanner *bufio.Scanner, data *string) {
	scanner.Scan()
	*data = scanner.Text()
}

func hexCharToBinaryString(c rune) string {
	switch c {
	case '0':
		return "0000"
	case '1':
		return "0001"
	case '2':

		return "0010"
	case '3':
		return "0011"
	case '4':
		return "0100"
	case '5':
		return "0101"
	case '6':
		return "0110"
	case '7':
		return "0111"
	case '8':
		return "1000"
	case '9':
		return "1001"
	case 'A':
		return "1010"
	case 'B':
		return "1011"
	case 'C':
		return "1100"
	case 'D':
		return "1101"
	case 'E':
		return "1110"
	case 'F':
		return "1111"
	default:
		return ""
	}
}

func parsePacket(binary string, index int) (uint64, int, uint64) {
	versionNumber, _ := strconv.ParseUint(binary[index:index+3], 2, 8)
	versionNumberTotal := versionNumber
	index += 3
	typeID, _ := strconv.ParseUint(binary[index:index+3], 2, 8)
	index += 3
	// fmt.Printf("v=%d, type=%d\n", versionNumber, typeID)

	var result uint64 = 0
	if typeID == literal {
		for {
			prefix := binary[index]
			index += 1
			part, _ := strconv.ParseUint(binary[index:index+4], 2, 8)
			index += 4
			result = result*16 + part
			if prefix == '0' {
				break
			}
		}
		// fmt.Printf("value: %d\n", value)
	} else {
		var subPacketLength uint64 = 0
		var numSubPackets uint64 = 0
		lengthTypeID := binary[index]
		index++
		if lengthTypeID == '0' {
			subPacketLength, _ = strconv.ParseUint(binary[index:index+15], 2, 16)
			index += 15
			// fmt.Printf("%d bits of subpackets\n", subPacketLength)
		} else {
			numSubPackets, _ = strconv.ParseUint(binary[index:index+11], 2, 16)
			index += 11
			// fmt.Printf("%d subpackets\n", numSubPackets)
		}

		startIndex := index
		numSubPacketsParsed := 0
		values := []uint64{}
		for {
			var subPacketVersionNumber uint64
			var value uint64
			subPacketVersionNumber, index, value = parsePacket(binary, index)
			values = append(values, value)
			versionNumberTotal += subPacketVersionNumber
			numSubPacketsParsed++
			if subPacketLength > 0 && index >= startIndex+int(subPacketLength) {
				break
			}
			if numSubPackets > 0 && numSubPacketsParsed >= int(numSubPackets) {
				break
			}
		}
		switch typeID {
		case sum:
			result = 0
			for _, v := range values {
				result += v
			}
		case product:
			result = 1
			for _, v := range values {
				result *= v
			}
		case minimum:
			result = values[0]
			for _, v := range values[1:] {
				if v < result {
					result = v
				}
			}
		case maximum:
			result = values[0]
			for _, v := range values[1:] {
				if v > result {
					result = v
				}
			}
		case greaterThan:
			if len(values) != 2 {
				panic("should be 2 subpackets!")
			}
			if values[0] > values[1] {
				result = 1
			} else {
				result = 0
			}
		case lessThan:
			if len(values) != 2 {
				panic("should be 2 subpackets!")
			}
			if values[0] < values[1] {
				result = 1
			} else {
				result = 0
			}
		case equalTo:
			if len(values) != 2 {
				panic("should be 2 subpackets!")
			}
			if values[0] == values[1] {
				result = 1
			} else {
				result = 0
			}
		}
	}
	return versionNumberTotal, index, result
}

func main() {
	filePath := flag.String("f", "", "file to read from")
	// isPart2 := flag.Bool("p", true, "part 2 switch")
	flag.Parse()

	var text string
	common.ParseFile(filePath, func(scanner *bufio.Scanner) {
		scanInput(scanner, &text)
	})

	var sb strings.Builder
	for _, c := range text {
		sb.WriteString(hexCharToBinaryString(c))
	}
	binary := sb.String()
	versionSum, _, result := parsePacket(binary, 0)
	fmt.Println("Version sum:", versionSum)
	fmt.Println("Result:", result)
}
