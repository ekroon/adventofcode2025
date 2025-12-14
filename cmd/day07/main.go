package main

import (
	"bufio"
	"fmt"
	"os"
)

type pos struct {
	row, col int
}

func part1(lines []string) int {
	grid := make([][]byte, len(lines))
	var start pos
	for r, line := range lines {
		grid[r] = []byte(line)
		for c, ch := range line {
			if ch == 'S' {
				start = pos{r, c}
			}
		}
	}

	rows, cols := len(grid), len(grid[0])

	// Track which splitters have been hit (by position)
	splitterHit := make(map[pos]bool)

	// BFS with beam positions (all beams move down)
	beams := []pos{start}
	visited := make(map[pos]bool)
	visited[start] = true

	for len(beams) > 0 {
		b := beams[0]
		beams = beams[1:]

		// Move down
		newPos := pos{b.row + 1, b.col}

		// Check bounds
		if newPos.row >= rows || newPos.col < 0 || newPos.col >= cols {
			continue
		}

		cell := grid[newPos.row][newPos.col]

		switch cell {
		case '.', 'S':
			// Continue down
			if !visited[newPos] {
				visited[newPos] = true
				beams = append(beams, newPos)
			}
		case '^':
			// Splitter: count it and split left and right (both continue down)
			splitterHit[newPos] = true
			// Left beam (col-1) and right beam (col+1), both continue down
			for _, dc := range []int{-1, 1} {
				splitPos := pos{newPos.row, newPos.col + dc}
				if splitPos.col >= 0 && splitPos.col < cols && !visited[splitPos] {
					visited[splitPos] = true
					beams = append(beams, splitPos)
				}
			}
		}
	}

	return len(splitterHit)
}

func part2(lines []string) int {
	grid := make([][]byte, len(lines))
	var start pos
	for r, line := range lines {
		grid[r] = []byte(line)
		for c, ch := range line {
			if ch == 'S' {
				start = pos{r, c}
			}
		}
	}

	rows, cols := len(grid), len(grid[0])

	// Count paths using dynamic programming
	// For each position, count how many ways to reach the bottom
	// We work bottom-up: memo[col] = number of paths from that column at current row

	// Start from the bottom row, work up
	// At each splitter, paths[col] = paths[col-1] + paths[col+1] from the row below

	// Actually, let's count paths going down from start
	// paths[col] = number of distinct paths reaching this column
	paths := make(map[int]int)
	paths[start.col] = 1

	for row := start.row; row < rows-1; row++ {
		nextPaths := make(map[int]int)

		for col, count := range paths {
			nextRow := row + 1
			if nextRow >= rows {
				continue
			}

			cell := grid[nextRow][col]

			switch cell {
			case '.', 'S':
				// Continue down, same number of paths
				nextPaths[col] += count
			case '^':
				// Splitter: choose left OR right, each choice is a separate path
				for _, dc := range []int{-1, 1} {
					newCol := col + dc
					if newCol >= 0 && newCol < cols {
						nextPaths[newCol] += count
					}
				}
			}
		}

		paths = nextPaths
	}

	// Sum all paths that reached the bottom
	total := 0
	for _, count := range paths {
		total += count
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
