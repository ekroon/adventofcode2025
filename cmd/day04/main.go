package main

import (
	"bufio"
	"fmt"
	"os"
)

func part1(lines []string) int {
	count := 0
	rows := len(lines)
	if rows == 0 {
		return 0
	}
	cols := len(lines[0])

	// 8 directions: N, NE, E, SE, S, SW, W, NW
	dirs := [][2]int{
		{-1, 0}, {-1, 1}, {0, 1}, {1, 1},
		{1, 0}, {1, -1}, {0, -1}, {-1, -1},
	}

	for r := range rows {
		for c := range cols {
			if lines[r][c] != '@' {
				continue
			}
			// Count adjacent stacks
			adjacent := 0
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols && lines[nr][nc] == '@' {
					adjacent++
				}
			}
			if adjacent < 4 {
				count++
			}
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
