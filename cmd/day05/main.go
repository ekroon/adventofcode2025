package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	low, high int
}

func parseInput(lines []string) ([]Range, []int) {
	var ranges []Range
	var numbers []int
	parsingRanges := true

	for _, line := range lines {
		if line == "" {
			parsingRanges = false
			continue
		}
		if parsingRanges {
			parts := strings.Split(line, "-")
			low, _ := strconv.Atoi(parts[0])
			high, _ := strconv.Atoi(parts[1])
			ranges = append(ranges, Range{low, high})
		} else {
			n, _ := strconv.Atoi(line)
			numbers = append(numbers, n)
		}
	}
	return ranges, numbers
}

func inAnyRange(n int, ranges []Range) bool {
	for _, r := range ranges {
		if n >= r.low && n <= r.high {
			return true
		}
	}
	return false
}

func part1(lines []string) int {
	ranges, numbers := parseInput(lines)
	count := 0
	for _, n := range numbers {
		if inAnyRange(n, ranges) {
			count++
		}
	}
	return count
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
