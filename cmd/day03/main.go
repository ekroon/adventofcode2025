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
