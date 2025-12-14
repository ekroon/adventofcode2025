package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// Shape represents a present shape as a list of (row, col) offsets
type Shape []struct{ r, c int }

// ShapeMask is a precomputed bitmask representation of a shape
type ShapeMask struct {
	rowMasks []uint64 // mask for each row (relative to top of shape)
	minRow   int      // minimum row offset
	maxRow   int      // maximum row offset
	maxCol   int      // maximum column offset
}

// Generate all unique orientations of a shape
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
	result := make(Shape, len(s))
	for i, p := range s {
		result[i] = struct{ r, c int }{p.c, -p.r}
	}
	return result
}

func flipShape(s Shape) Shape {
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
	points := make([]string, len(s))
	for i, p := range s {
		points[i] = fmt.Sprintf("%d,%d", p.r, p.c)
	}
	slices.Sort(points)
	return strings.Join(points, ";")
}

// Convert shape to bitmask form
func shapeToMask(s Shape) ShapeMask {
	if len(s) == 0 {
		return ShapeMask{}
	}

	maxR, maxC := 0, 0
	for _, p := range s {
		maxR = max(maxR, p.r)
		maxC = max(maxC, p.c)
	}

	rowMasks := make([]uint64, maxR+1)
	for _, p := range s {
		rowMasks[p.r] |= 1 << p.c
	}

	return ShapeMask{
		rowMasks: rowMasks,
		minRow:   0,
		maxRow:   maxR,
		maxCol:   maxC,
	}
}

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
	counts        []int
}

func isRegionLine(line string) bool {
	return strings.Contains(line, "x") && strings.Contains(line, ": ")
}

func parseInput(lines []string) ([][]ShapeMask, []Region, int) {
	regionStart := -1
	for i, line := range lines {
		if isRegionLine(line) {
			regionStart = i
			break
		}
	}

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

	// Convert to masks for all orientations
	allMasks := make([][]ShapeMask, len(baseShapes))
	cellCount := 0
	for i, s := range baseShapes {
		orientations := allOrientations(s)
		allMasks[i] = make([]ShapeMask, len(orientations))
		for j, o := range orientations {
			allMasks[i][j] = shapeToMask(o)
		}
		if len(s) > 0 {
			cellCount = len(s)
		}
	}

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

	return allMasks, regions, cellCount
}

// Grid for solving using row bitmasks
type Grid struct {
	width, height int
	rows          []uint64
}

func newGrid(w, h int) *Grid {
	return &Grid{
		width:  w,
		height: h,
		rows:   make([]uint64, h),
	}
}

func (g *Grid) canPlace(m *ShapeMask, startR, startC int) bool {
	if startR < 0 || startR+m.maxRow >= g.height || startC < 0 || startC+m.maxCol >= g.width {
		return false
	}
	for i, mask := range m.rowMasks {
		shiftedMask := mask << startC
		if g.rows[startR+i]&shiftedMask != 0 {
			return false
		}
	}
	return true
}

func (g *Grid) place(m *ShapeMask, startR, startC int) {
	for i, mask := range m.rowMasks {
		g.rows[startR+i] |= mask << startC
	}
}

func (g *Grid) remove(m *ShapeMask, startR, startC int) {
	for i, mask := range m.rowMasks {
		g.rows[startR+i] &^= mask << startC
	}
}

func (g *Grid) firstEmpty() (int, int) {
	for r := 0; r < g.height; r++ {
		invRow := ^g.rows[r]
		if invRow != 0 {
			c := bits.TrailingZeros64(invRow)
			if c < g.width {
				return r, c
			}
		}
	}
	return -1, -1
}

func (g *Grid) setCell(r, c int) {
	g.rows[r] |= 1 << c
}

func (g *Grid) clearCell(r, c int) {
	g.rows[r] &^= 1 << c
}

// ShapeEntry for the shape list
type ShapeEntry struct {
	shapeIdx  int
	cellCount int
}

// Solver for backtracking
type Solver struct {
	grid     *Grid
	allMasks [][]ShapeMask
	shapes   []ShapeEntry
}

func (s *Solver) solve(idx int, skipsLeft int) bool {
	if idx == len(s.shapes) {
		return true
	}

	r, c := s.grid.firstEmpty()
	if r < 0 {
		return false
	}

	for i := idx; i < len(s.shapes); i++ {
		if i > idx && s.shapes[i].shapeIdx == s.shapes[idx].shapeIdx {
			continue
		}

		s.shapes[idx], s.shapes[i] = s.shapes[i], s.shapes[idx]
		shapeIdx := s.shapes[idx].shapeIdx

		for mi := range s.allMasks[shapeIdx] {
			m := &s.allMasks[shapeIdx][mi]
			for dr := 0; dr <= m.maxRow; dr++ {
				mask := m.rowMasks[dr]
				for dc := 0; dc <= m.maxCol; dc++ {
					if mask&(1<<dc) != 0 {
						startR := r - dr
						startC := c - dc
						if s.grid.canPlace(m, startR, startC) {
							s.grid.place(m, startR, startC)
							if s.solve(idx+1, skipsLeft) {
								return true
							}
							s.grid.remove(m, startR, startC)
						}
					}
				}
			}
		}

		s.shapes[idx], s.shapes[i] = s.shapes[i], s.shapes[idx]
	}

	if skipsLeft > 0 {
		s.grid.setCell(r, c)
		result := s.solve(idx, skipsLeft-1)
		s.grid.clearCell(r, c)
		return result
	}

	return false
}

// Greedy placement - try row by row, packing tightly
func greedyPlace(g *Grid, allMasks [][]ShapeMask, shapes []ShapeEntry) bool {
	for _, se := range shapes {
		placed := false
		// Try each orientation
		for mi := range allMasks[se.shapeIdx] {
			if placed {
				break
			}
			m := &allMasks[se.shapeIdx][mi]
			maxStartR := g.height - m.maxRow - 1
			maxStartC := g.width - m.maxCol - 1
			// Try positions, prioritizing top-left
			for startR := 0; startR <= maxStartR && !placed; startR++ {
				for startC := 0; startC <= maxStartC && !placed; startC++ {
					if g.canPlace(m, startR, startC) {
						g.place(m, startR, startC)
						placed = true
					}
				}
			}
		}
		if !placed {
			return false
		}
	}
	return true
}

func canFit(allMasks [][]ShapeMask, region Region, cellCount int) bool {
	// Build shape list
	var shapes []ShapeEntry
	totalCells := 0
	for shapeIdx, count := range region.counts {
		if count > 0 && len(allMasks[shapeIdx]) > 0 {
			for range count {
				shapes = append(shapes, ShapeEntry{shapeIdx: shapeIdx, cellCount: cellCount})
			}
			totalCells += count * cellCount
		}
	}

	if len(shapes) == 0 {
		return true
	}

	gridArea := region.width * region.height
	if totalCells > gridArea {
		return false
	}

	// Try greedy first
	g := newGrid(region.width, region.height)
	if greedyPlace(g, allMasks, shapes) {
		return true
	}

	// Fall back to backtracking
	solver := &Solver{
		grid:     newGrid(region.width, region.height),
		allMasks: allMasks,
		shapes:   shapes,
	}

	skipsAllowed := gridArea - totalCells
	return solver.solve(0, skipsAllowed)
}

func part1(lines []string) int {
	allMasks, regions, cellCount := parseInput(lines)

	var wg sync.WaitGroup
	var count atomic.Int64

	for _, region := range regions {
		wg.Add(1)
		go func(r Region) {
			defer wg.Done()
			if canFit(allMasks, r, cellCount) {
				count.Add(1)
			}
		}(region)
	}

	wg.Wait()
	return int(count.Load())
}

func part2(lines []string) int {
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
