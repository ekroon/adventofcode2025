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
	if len(lines) < 2 {
		return 0
	}

	// Find max width
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Pad all lines to same width
	padded := make([]string, len(lines))
	for i, line := range lines {
		padded[i] = line + strings.Repeat(" ", maxWidth-len(line))
	}

	numRows := len(padded)
	total := 0
	var currentNumbers []int

	// Read from right to left, one column at a time
	for col := maxWidth - 1; col >= 0; col-- {
		// Check if operation row has an operator
		opChar := padded[numRows-1][col]

		// Build ONE number from this column by reading digits top-to-bottom
		num := 0
		hasDigit := false
		for row := 0; row < numRows-1; row++ {
			ch := padded[row][col]
			if ch >= '0' && ch <= '9' {
				num = num*10 + int(ch-'0')
				hasDigit = true
			}
		}

		if hasDigit {
			currentNumbers = append(currentNumbers, num)
		}

		if opChar == '+' || opChar == '*' {
			// Apply operation to all collected numbers
			if len(currentNumbers) > 0 {
				result := currentNumbers[0]
				for _, n := range currentNumbers[1:] {
					if opChar == '+' {
						result += n
					} else {
						result *= n
					}
				}
				total += result
			}
			currentNumbers = nil
		}
	}

	return total
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
