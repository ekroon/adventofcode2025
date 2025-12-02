package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// isInvalidDouble checks if a number consists of a sequence of digits repeated exactly twice
// e.g., 11 (1+1), 123123 (123+123) are invalid
func isInvalidDouble(id int) bool {
	s := strconv.Itoa(id)
	// Must have even length to be a doubled sequence
	if len(s)%2 != 0 {
		return false
	}
	half := len(s) / 2
	return s[:half] == s[half:]
}

// isInvalidRepeated checks if a number consists of a sequence repeated 2+ times
// e.g., 11, 111, 1111, 123123, 123123123 are all invalid
func isInvalidRepeated(id int) bool {
	s := strconv.Itoa(id)
	// Try each possible pattern length (1 to len/2)
	for patLen := 1; patLen <= len(s)/2; patLen++ {
		if len(s)%patLen != 0 {
			continue
		}
		pattern := s[:patLen]
		repeated := strings.Repeat(pattern, len(s)/patLen)
		if repeated == s {
			return true
		}
	}
	return false
}

func part1(lines []string) int {
	if len(lines) == 0 {
		return 0
	}

	sum := 0
	// Parse ranges from the line: "11-22,1234-2345"
	for r := range strings.SplitSeq(lines[0], ",") {
		startStr, endStr, _ := strings.Cut(r, "-")
		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)

		for id := start; id <= end; id++ {
			if isInvalidDouble(id) {
				sum += id
			}
		}
	}
	return sum
}

func part2(lines []string) int {
	if len(lines) == 0 {
		return 0
	}

	sum := 0
	for r := range strings.SplitSeq(lines[0], ",") {
		startStr, endStr, _ := strings.Cut(r, "-")
		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)

		for id := start; id <= end; id++ {
			if isInvalidRepeated(id) {
				sum += id
			}
		}
	}
	return sum
}

func main() {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println("Part 1:", part1(lines))
	fmt.Println("Part 2:", part2(lines))
}
