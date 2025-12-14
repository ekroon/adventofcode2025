package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x, y int
}

func parseInput(lines []string) []point {
	points := make([]point, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		points = append(points, point{x, y})
	}
	return points
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func part1(lines []string) int {
	points := parseInput(lines)
	maxArea := 0

	// Check all pairs of points as opposite corners
	for i := range len(points) {
		for j := i + 1; j < len(points); j++ {
			p1, p2 := points[i], points[j]
			// For opposite corners, x and y must both differ
			if p1.x != p2.x && p1.y != p2.y {
				area := (abs(p2.x-p1.x) + 1) * (abs(p2.y-p1.y) + 1)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	return maxArea
}

func part2(lines []string) int {
	points := parseInput(lines)
	if len(points) < 2 {
		return 0
	}

	// Build set of red tiles
	redSet := make(map[point]bool)
	for _, p := range points {
		redSet[p] = true
	}

	// Build set of green tiles (edges between consecutive red tiles)
	greenSet := make(map[point]bool)
	for i := range len(points) {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]

		// Draw line between p1 and p2 (they share either x or y)
		if p1.x == p2.x {
			// Vertical line
			minY, maxY := min(p1.y, p2.y), max(p1.y, p2.y)
			for y := minY; y <= maxY; y++ {
				pt := point{p1.x, y}
				if !redSet[pt] {
					greenSet[pt] = true
				}
			}
		} else {
			// Horizontal line
			minX, maxX := min(p1.x, p2.x), max(p1.x, p2.x)
			for x := minX; x <= maxX; x++ {
				pt := point{x, p1.y}
				if !redSet[pt] {
					greenSet[pt] = true
				}
			}
		}
	}

	// Find bounding box
	minX, maxX := points[0].x, points[0].x
	minY, maxY := points[0].y, points[0].y
	for _, p := range points {
		minX, maxX = min(minX, p.x), max(maxX, p.x)
		minY, maxY = min(minY, p.y), max(maxY, p.y)
	}

	// Combine red and green as "loop" tiles, then flood fill interior
	loopSet := make(map[point]bool)
	for p := range redSet {
		loopSet[p] = true
	}
	for p := range greenSet {
		loopSet[p] = true
	}

	// Flood fill from outside to find exterior tiles
	// Expand bounding box by 1 to ensure we can reach around the loop
	exterior := make(map[point]bool)
	queue := []point{{minX - 1, minY - 1}}
	exterior[queue[0]] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, d := range []point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			np := point{cur.x + d.x, cur.y + d.y}
			if np.x < minX-1 || np.x > maxX+1 || np.y < minY-1 || np.y > maxY+1 {
				continue
			}
			if exterior[np] || loopSet[np] {
				continue
			}
			exterior[np] = true
			queue = append(queue, np)
		}
	}

	// Interior tiles are those not on loop and not exterior
	// Add interior to green set
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			p := point{x, y}
			if !loopSet[p] && !exterior[p] {
				greenSet[p] = true
			}
		}
	}

	// Valid tiles = red + green
	validSet := make(map[point]bool)
	for p := range redSet {
		validSet[p] = true
	}
	for p := range greenSet {
		validSet[p] = true
	}

	// Check all pairs of red points as opposite corners
	maxArea := 0
	for i := range len(points) {
		for j := i + 1; j < len(points); j++ {
			p1, p2 := points[i], points[j]
			if p1.x == p2.x || p1.y == p2.y {
				continue
			}

			// Check if entire rectangle is valid
			x1, x2 := min(p1.x, p2.x), max(p1.x, p2.x)
			y1, y2 := min(p1.y, p2.y), max(p1.y, p2.y)

			valid := true
			for x := x1; x <= x2 && valid; x++ {
				for y := y1; y <= y2 && valid; y++ {
					if !validSet[point{x, y}] {
						valid = false
					}
				}
			}

			if valid {
				area := (x2 - x1 + 1) * (y2 - y1 + 1)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}

	return maxArea
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
