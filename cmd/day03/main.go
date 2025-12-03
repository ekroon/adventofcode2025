package main

import (
	"bufio"
	"fmt"
	"os"
)

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		// Find the max two-digit number by picking two digits in order
		maxNum := 0
		for i := 0; i < len(line)-1; i++ {
			for j := i + 1; j < len(line); j++ {
				num := int(line[i]-'0')*10 + int(line[j]-'0')
				maxNum = max(maxNum, num)
			}
		}
		total += maxNum
	}
	return total
}

func part2(lines []string) int {
	total := 0
	for _, line := range lines {
		total += maxNumber(line, 12)
	}
	return total
}

// maxNumber finds the maximum number by picking `count` digits from line in order
func maxNumber(line string, count int) int {
	n := len(line)
	if count > n {
		count = n
	}

	// Greedy approach: for each position, pick the largest digit that still
	// leaves enough digits remaining for the rest of the number
	result := 0
	start := 0
	for i := 0; i < count; i++ {
		// We need (count - i - 1) more digits after this one
		// So we can pick from start to n-(count-i-1)-1 = n-count+i
		end := n - count + i
		bestDigit := byte('0')
		bestIdx := start
		for j := start; j <= end; j++ {
			if line[j] > bestDigit {
				bestDigit = line[j]
				bestIdx = j
			}
		}
		result = result*10 + int(bestDigit-'0')
		start = bestIdx + 1
	}
	return result
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
