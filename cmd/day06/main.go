package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func part1(lines []string) int {
	if len(lines) < 2 {
		return 0
	}

	// Parse number rows (all but last line)
	var columns [][]int
	for _, line := range lines[:len(lines)-1] {
		fields := strings.Fields(line)
		for i, f := range fields {
			num, _ := strconv.Atoi(f)
			if i >= len(columns) {
				columns = append(columns, []int{})
			}
			columns[i] = append(columns[i], num)
		}
	}

	// Parse operations (last line)
	ops := strings.Fields(lines[len(lines)-1])

	// Calculate each column with its operation
	total := 0
	for i, col := range columns {
		op := ops[i]
		result := col[0]
		for _, num := range col[1:] {
			if op == "*" {
				result *= num
			} else {
				result += num
			}
		}
		total += result
	}

	return total
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
