package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// isInvalid checks if a number consists of a sequence of digits repeated twice
// e.g., 11 (1+1), 123123 (123+123) are invalid
func isInvalid(id int) bool {
	s := strconv.Itoa(id)
	// Must have even length to be a doubled sequence
	if len(s)%2 != 0 {
		return false
	}
	half := len(s) / 2
	return s[:half] == s[half:]
}

func part1(lines []string) int {
	if len(lines) == 0 {
		return 0
	}

	sum := 0
	// Parse ranges from the line: "11-22,1234-2345"
	ranges := strings.Split(lines[0], ",")
	for _, r := range ranges {
		parts := strings.Split(r, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])

		for id := start; id <= end; id++ {
			if isInvalid(id) {
				sum += id
			}
		}
	}
	return sum
}

func part2(lines []string) int {
	// TODO: implement
	return 0
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
