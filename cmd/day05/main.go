package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

func mergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return nil
	}

	// Sort by low bound
	sorted := make([]Range, len(ranges))
	copy(sorted, ranges)
	slices.SortFunc(sorted, func(a, b Range) int {
		return a.low - b.low
	})

	merged := []Range{sorted[0]}
	for _, r := range sorted[1:] {
		last := &merged[len(merged)-1]
		if r.low <= last.high+1 {
			// Overlapping or adjacent, extend the range
			last.high = max(last.high, r.high)
		} else {
			// Gap, start new range
			merged = append(merged, r)
		}
	}
	return merged
}

func part2(lines []string) int {
	ranges, _ := parseInput(lines)
	merged := mergeRanges(ranges)

	count := 0
	for _, r := range merged {
		count += r.high - r.low + 1
	}
	return count
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
