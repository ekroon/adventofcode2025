package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

// segment represents an axis-aligned line segment
type segment struct {
	x1, y1, x2, y2 int
	horizontal     bool
}

func part2(lines []string) int {
	points := parseInput(lines)
	if len(points) < 3 {
		return 0
	}

	// Build polygon segments from consecutive red tiles
	segments := make([]segment, len(points))
	for i := range points {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		seg := segment{
			x1:         min(p1.x, p2.x),
			y1:         min(p1.y, p2.y),
			x2:         max(p1.x, p2.x),
			y2:         max(p1.y, p2.y),
			horizontal: p1.y == p2.y,
		}
		segments[i] = seg
	}

	// Collect all unique X and Y coordinates
	xSet := make(map[int]bool)
	ySet := make(map[int]bool)
	for _, p := range points {
		xSet[p.x] = true
		ySet[p.y] = true
	}

	// Convert to sorted slices
	xCoords := make([]int, 0, len(xSet))
	for x := range xSet {
		xCoords = append(xCoords, x)
	}
	yCoords := make([]int, 0, len(ySet))
	for y := range ySet {
		yCoords = append(yCoords, y)
	}
	slices.Sort(xCoords)
	slices.Sort(yCoords)

	// Create reverse maps for O(1) lookup
	xIndex := make(map[int]int)
	for i, x := range xCoords {
		xIndex[x] = i
	}
	yIndex := make(map[int]int)
	for i, y := range yCoords {
		yIndex[y] = i
	}

	// Build compressed grid: each cell represents the region between consecutive coordinates
	// A cell at (i,j) in compressed space represents the rectangle from
	// (xCoords[i], yCoords[j]) to (xCoords[i+1]-1, yCoords[j+1]-1) in real space
	// But for our purposes, we track whether entire cells are inside the polygon

	compW := len(xCoords)
	compH := len(yCoords)

	// For each compressed cell, determine if it's inside the polygon
	// We use ray casting: a point is inside if the number of edge crossings
	// to the left (going to x=-inf) is odd

	// isInside checks if a point is strictly inside the polygon using ray casting
	// We cast a ray from (x, y) to (x, -infinity) and count crossings
	isInsidePolygon := func(x, y int) bool {
		crossings := 0
		for _, seg := range segments {
			if seg.horizontal {
				// Horizontal segment at seg.y1
				// Ray going down from (x,y) crosses if:
				// - seg.y1 < y (segment is below our point)
				// - x is within the segment's x range (exclusive on one end for consistency)
				if seg.y1 < y && seg.x1 <= x && x < seg.x2 {
					crossings++
				}
			}
			// Vertical segments don't cross a vertical ray
		}
		return crossings%2 == 1
	}

	// isOnBoundary checks if a point is on the polygon boundary
	isOnBoundary := func(x, y int) bool {
		for _, seg := range segments {
			if seg.horizontal {
				if y == seg.y1 && seg.x1 <= x && x <= seg.x2 {
					return true
				}
			} else {
				if x == seg.x1 && seg.y1 <= y && y <= seg.y2 {
					return true
				}
			}
		}
		return false
	}

	// isValidPoint checks if a point is inside or on the boundary
	isValidPoint := func(x, y int) bool {
		return isOnBoundary(x, y) || isInsidePolygon(x, y)
	}

	// Build compressed grid of "inside" status
	// inside[i][j] = true if the entire cell from (xCoords[i], yCoords[j]) to
	// (xCoords[i+1]-1, yCoords[j+1]-1) is inside (or on boundary of) the polygon
	// For corner cells (i = compW-1 or j = compH-1), they're just the single points
	inside := make([][]bool, compW)
	for i := range inside {
		inside[i] = make([]bool, compH)
	}

	for i := 0; i < compW; i++ {
		for j := 0; j < compH; j++ {
			// Check all 4 corners of this compressed cell
			x1 := xCoords[i]
			y1 := yCoords[j]
			x2 := x1
			y2 := y1
			if i+1 < compW {
				x2 = xCoords[i+1] - 1
			}
			if j+1 < compH {
				y2 = yCoords[j+1] - 1
			}

			// A cell is fully inside if all 4 corners are valid
			// and no polygon edge crosses through the cell's interior
			if !isValidPoint(x1, y1) || !isValidPoint(x2, y1) ||
				!isValidPoint(x1, y2) || !isValidPoint(x2, y2) {
				continue
			}

			// Check if any polygon edge crosses through this cell's interior
			// An edge crosses if it's strictly inside the cell bounds (not on boundary)
			edgeCrosses := false
			for _, seg := range segments {
				if seg.horizontal {
					// Horizontal edge at y=seg.y1 from x=seg.x1 to x=seg.x2
					// Crosses cell interior if seg.y1 is strictly between y1 and y2
					// and segment overlaps with cell's x range
					if y1 < seg.y1 && seg.y1 < y2 {
						if seg.x1 < x2 && seg.x2 > x1 {
							edgeCrosses = true
							break
						}
					}
				} else {
					// Vertical edge at x=seg.x1 from y=seg.y1 to y=seg.y2
					if x1 < seg.x1 && seg.x1 < x2 {
						if seg.y1 < y2 && seg.y2 > y1 {
							edgeCrosses = true
							break
						}
					}
				}
			}

			inside[i][j] = !edgeCrosses
		}
	}

	// Build prefix sum over compressed grid for O(1) range queries
	// prefix[i][j] = number of inside cells in [0..i-1][0..j-1]
	prefix := make([][]int, compW+1)
	for i := range prefix {
		prefix[i] = make([]int, compH+1)
	}

	for i := 0; i < compW; i++ {
		for j := 0; j < compH; j++ {
			val := 0
			if inside[i][j] {
				val = 1
			}
			prefix[i+1][j+1] = val + prefix[i][j+1] + prefix[i+1][j] - prefix[i][j]
		}
	}

	// countInsideCells returns number of inside cells in compressed range [ci1..ci2-1][cj1..cj2-1]
	countInsideCells := func(ci1, cj1, ci2, cj2 int) int {
		return prefix[ci2][cj2] - prefix[ci1][cj2] - prefix[ci2][cj1] + prefix[ci1][cj1]
	}

	// For each pair of red points, check if the rectangle between them is fully valid
	maxArea := 0
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			p1, p2 := points[i], points[j]
			if p1.x == p2.x || p1.y == p2.y {
				continue
			}

			x1, x2 := min(p1.x, p2.x), max(p1.x, p2.x)
			y1, y2 := min(p1.y, p2.y), max(p1.y, p2.y)

			// Get compressed indices
			ci1, ci2 := xIndex[x1], xIndex[x2]
			cj1, cj2 := yIndex[y1], yIndex[y2]

			// Total compressed cells in this range
			totalCells := (ci2 - ci1) * (cj2 - cj1)

			// Count how many are inside
			insideCount := countInsideCells(ci1, cj1, ci2, cj2)

			// Rectangle is valid if all cells are inside
			if insideCount == totalCells {
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
