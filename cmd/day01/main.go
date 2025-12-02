package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func part1(lines []string) int {
	dial := 50
	password := 0

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		direction := line[0]
		amount, _ := strconv.Atoi(line[1:])

		if direction == 'L' {
			dial -= amount
		} else if direction == 'R' {
			dial += amount
		}

		// Wrap dial to 0-99 range
		dial = ((dial % 100) + 100) % 100

		if dial == 0 {
			password++
		}
	}

	return password
}

func part2(lines []string) int {
	dial := 50
	password := 0

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		direction := line[0]
		amount, _ := strconv.Atoi(line[1:])

		for i := 0; i < amount; i++ {
			if direction == 'L' {
				dial--
				if dial < 0 {
					dial = 99
				}
			} else if direction == 'R' {
				dial++
				if dial > 99 {
					dial = 0
				}
			}

			if dial == 0 {
				password++
			}
		}
	}

	return password
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
