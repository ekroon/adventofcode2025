package main

import (
	"bufio"
	"fmt"
	"os"
)

func part1(lines []string) int {
	// TODO: implement
	return 0
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
