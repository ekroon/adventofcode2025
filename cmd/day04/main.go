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
	rows := len(lines)
	if rows == 0 {
		return 0
	}
	cols := len(lines[0])

	// Convert to mutable grid
	grid := make([][]byte, rows)
	for r := range rows {
		grid[r] = []byte(lines[r])
	}

	// 8 directions: N, NE, E, SE, S, SW, W, NW
	dirs := [][2]int{
		{-1, 0}, {-1, 1}, {0, 1}, {1, 1},
		{1, 0}, {1, -1}, {0, -1}, {-1, -1},
	}

	totalRemoved := 0

	for {
		// Find all stacks to remove this round
		var toRemove [][2]int
		for r := range rows {
			for c := range cols {
				if grid[r][c] != '@' {
					continue
				}
				adjacent := 0
				for _, d := range dirs {
					nr, nc := r+d[0], c+d[1]
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '@' {
						adjacent++
					}
				}
				if adjacent < 4 {
					toRemove = append(toRemove, [2]int{r, c})
				}
			}
		}

		if len(toRemove) == 0 {
			break
		}

		// Remove all marked stacks
		for _, pos := range toRemove {
			grid[pos[0]][pos[1]] = '.'
		}
		totalRemoved += len(toRemove)
	}

	return totalRemoved
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
