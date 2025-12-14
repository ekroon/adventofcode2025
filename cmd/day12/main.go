package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Shape represents a present shape as a list of (row, col) offsets
type Shape []struct{ r, c int }

// Generate all 8 orientations (4 rotations x 2 flips) of a shape
func allOrientations(s Shape) []Shape {
	orientations := make(map[string]Shape)

	current := s
	for flip := 0; flip < 2; flip++ {
		for rot := 0; rot < 4; rot++ {
			normalized := normalize(current)
			key := shapeKey(normalized)
			if _, exists := orientations[key]; !exists {
				orientations[key] = normalized
			}
			current = rotate90(current)
		}
		current = flipShape(s)
	}

	result := make([]Shape, 0, len(orientations))
	for _, shape := range orientations {
		result = append(result, shape)
	}
	return result
}

func rotate90(s Shape) Shape {
	// Rotate 90 degrees clockwise: (r, c) -> (c, -r)
	result := make(Shape, len(s))
	for i, p := range s {
		result[i] = struct{ r, c int }{p.c, -p.r}
	}
	return result
}

func flipShape(s Shape) Shape {
	// Flip horizontally: (r, c) -> (r, -c)
	result := make(Shape, len(s))
	for i, p := range s {
		result[i] = struct{ r, c int }{p.r, -p.c}
	}
	return result
}

func normalize(s Shape) Shape {
	if len(s) == 0 {
		return s
	}
	minR, minC := s[0].r, s[0].c
	for _, p := range s {
		minR = min(minR, p.r)
		minC = min(minC, p.c)
	}
	result := make(Shape, len(s))
	for i, p := range s {
		result[i] = struct{ r, c int }{p.r - minR, p.c - minC}
	}
	return result
}

func shapeKey(s Shape) string {
	// Create a canonical string representation
	points := make([]string, len(s))
	for i, p := range s {
		points[i] = fmt.Sprintf("%d,%d", p.r, p.c)
	}
	// Sort for consistency
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			if points[i] > points[j] {
				points[i], points[j] = points[j], points[i]
			}
		}
	}
	return strings.Join(points, ";")
}

// Grid represents the region being filled
type Grid struct {
	width, height int
	cells         [][]bool // true if occupied
}

func newGrid(width, height int) *Grid {
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return &Grid{width: width, height: height, cells: cells}
}

func (g *Grid) canPlace(s Shape, startR, startC int) bool {
	for _, p := range s {
		r, c := startR+p.r, startC+p.c
		if r < 0 || r >= g.height || c < 0 || c >= g.width {
			return false
		}
		if g.cells[r][c] {
			return false
		}
	}
	return true
}

func (g *Grid) place(s Shape, startR, startC int) {
	for _, p := range s {
		g.cells[startR+p.r][startC+p.c] = true
	}
}

func (g *Grid) remove(s Shape, startR, startC int) {
	for _, p := range s {
		g.cells[startR+p.r][startC+p.c] = false
	}
}

// Find the first empty cell (top-left to bottom-right)
func (g *Grid) firstEmpty() (int, int, bool) {
	for r := 0; r < g.height; r++ {
		for c := 0; c < g.width; c++ {
			if !g.cells[r][c] {
				return r, c, true
			}
		}
	}
	return -1, -1, false
}

// Parse a shape from lines like "###", "##.", "##."
func parseShape(lines []string) Shape {
	var shape Shape
	for r, line := range lines {
		for c, ch := range line {
			if ch == '#' {
				shape = append(shape, struct{ r, c int }{r, c})
			}
		}
	}
	return shape
}

// Region problem
type Region struct {
	width, height int
	counts        []int // count of each shape needed
}

// isRegionLine checks if a line is a region definition (WxH: ...)
func isRegionLine(line string) bool {
	return strings.Contains(line, "x") && strings.Contains(line, ": ")
}

func parseInput(lines []string) ([][]Shape, []Region) {
	// Find first region line - separates shapes from regions
	regionStart := -1
	for i, line := range lines {
		if isRegionLine(line) {
			regionStart = i
			break
		}
	}

	// Parse shapes
	shapeLines := lines[:regionStart]
	baseShapes := []Shape{}
	var currentShapeLines []string
	for _, line := range shapeLines {
		if strings.Contains(line, ":") {
			if len(currentShapeLines) > 0 {
				baseShapes = append(baseShapes, parseShape(currentShapeLines))
			}
			currentShapeLines = nil
		} else if line != "" {
			currentShapeLines = append(currentShapeLines, line)
		}
	}
	if len(currentShapeLines) > 0 {
		baseShapes = append(baseShapes, parseShape(currentShapeLines))
	}

	// Precompute all orientations for each shape
	allShapes := make([][]Shape, len(baseShapes))
	for i, s := range baseShapes {
		allShapes[i] = allOrientations(s)
	}

	// Parse regions
	regionLines := lines[regionStart:]
	regions := []Region{}
	for _, line := range regionLines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		dims := strings.Split(parts[0], "x")
		width, _ := strconv.Atoi(dims[0])
		height, _ := strconv.Atoi(dims[1])

		countStrs := strings.Fields(parts[1])
		counts := make([]int, len(countStrs))
		for i, cs := range countStrs {
			counts[i], _ = strconv.Atoi(cs)
		}

		regions = append(regions, Region{width: width, height: height, counts: counts})
	}

	return allShapes, regions
}

// Build a flat list of shape indices to place (e.g., if counts = [2, 0, 1], result = [0, 0, 2])
func buildShapeList(counts []int) []int {
	var list []int
	for shapeIdx, count := range counts {
		for i := 0; i < count; i++ {
			list = append(list, shapeIdx)
		}
	}
	return list
}

// Solve tries to place all shapes in the list using backtracking
// skipsLeft is the number of cells we can still leave empty
func solve(g *Grid, allShapes [][]Shape, shapeList []int, idx int, skipsLeft int) bool {
	if idx == len(shapeList) {
		return true // All shapes placed
	}

	// Find first empty cell - we try to place a shape covering this cell
	r, c, found := g.firstEmpty()
	if !found {
		// Grid is full but we still have shapes to place
		return false
	}

	// Try each remaining shape in the list (swap-based approach)
	for i := idx; i < len(shapeList); i++ {
		// Skip if same shape type as already tried at this position
		if i > idx && shapeList[i] == shapeList[idx] {
			continue
		}

		shapeList[idx], shapeList[i] = shapeList[i], shapeList[idx]
		shapeIdx := shapeList[idx]

		// Try each orientation of this shape
		for _, orientation := range allShapes[shapeIdx] {
			// Try placing so that the shape covers cell (r, c)
			for _, offset := range orientation {
				startR := r - offset.r
				startC := c - offset.c
				if g.canPlace(orientation, startR, startC) {
					g.place(orientation, startR, startC)
					if solve(g, allShapes, shapeList, idx+1, skipsLeft) {
						return true
					}
					g.remove(orientation, startR, startC)
				}
			}
		}

		shapeList[idx], shapeList[i] = shapeList[i], shapeList[idx]
	}

	// If no shape could cover this cell and we have skips left, try skipping
	if skipsLeft > 0 {
		g.cells[r][c] = true // temporarily block this cell
		result := solve(g, allShapes, shapeList, idx, skipsLeft-1)
		g.cells[r][c] = false // unblock
		return result
	}

	return false
}

func canFit(allShapes [][]Shape, region Region) bool {
	g := newGrid(region.width, region.height)
	shapeList := buildShapeList(region.counts)

	if len(shapeList) == 0 {
		return true
	}

	// Calculate total cells needed vs available
	totalCells := 0
	for i, count := range region.counts {
		if count > 0 && len(allShapes[i]) > 0 {
			totalCells += count * len(allShapes[i][0])
		}
	}
	gridArea := region.width * region.height
	if totalCells > gridArea {
		return false
	}

	skipsAllowed := gridArea - totalCells
	return solve(g, allShapes, shapeList, 0, skipsAllowed)
}

func part1(lines []string) int {
	allShapes, regions := parseInput(lines)

	count := 0
	for _, region := range regions {
		if canFit(allShapes, region) {
			count++
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
