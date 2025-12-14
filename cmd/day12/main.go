package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"runtime"
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
	rowMasks  []uint64 // mask for each row (relative to top of shape)
	bitShifts []int    // precomputed bit positions for blocked mask
	minRow    int      // minimum row offset
	maxRow    int      // maximum row offset
	maxCol    int      // maximum column offset
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

	// Precompute bit positions for blocked mask calculation
	// Each entry is (rowIndex << 8) | bitPos for efficient iteration
	var bitShifts []int
	for i, mask := range rowMasks {
		for bit := mask; bit != 0; {
			bitPos := bits.TrailingZeros64(bit)
			bitShifts = append(bitShifts, (i<<8)|bitPos)
			bit &= bit - 1
		}
	}

	return ShapeMask{
		rowMasks:  rowMasks,
		bitShifts: bitShifts,
		minRow:    0,
		maxRow:    maxR,
		maxCol:    maxC,
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
// Uses bitmask operations to find valid positions faster
func greedyPlace(g *Grid, allMasks [][]ShapeMask, shapes []ShapeEntry) bool {
	for _, se := range shapes {
		placed := false
		orientations := allMasks[se.shapeIdx]

		// Try each orientation
	orientLoop:
		for mi := range orientations {
			m := &orientations[mi]
			maxStartR := g.height - m.maxRow - 1
			maxStartC := g.width - m.maxCol - 1
			if maxStartR < 0 || maxStartC < 0 {
				continue
			}
			rowMasks := m.rowMasks
			nRows := len(rowMasks)
			validMask := uint64((1 << (maxStartC + 1)) - 1)
			shifts := m.bitShifts
			rows := g.rows
			nShifts := len(shifts)

			// Try positions row by row
			for startR := 0; startR <= maxStartR; startR++ {
				// Compute combined blocked mask using precomputed bit shifts
				// Unroll for common cell counts (5, 6, 7)
				var blocked uint64
				switch nShifts {
				case 5:
					s0, s1, s2, s3, s4 := shifts[0], shifts[1], shifts[2], shifts[3], shifts[4]
					blocked = rows[startR+(s0>>8)]>>(s0&0xFF) |
						rows[startR+(s1>>8)]>>(s1&0xFF) |
						rows[startR+(s2>>8)]>>(s2&0xFF) |
						rows[startR+(s3>>8)]>>(s3&0xFF) |
						rows[startR+(s4>>8)]>>(s4&0xFF)
				case 6:
					s0, s1, s2, s3, s4, s5 := shifts[0], shifts[1], shifts[2], shifts[3], shifts[4], shifts[5]
					blocked = rows[startR+(s0>>8)]>>(s0&0xFF) |
						rows[startR+(s1>>8)]>>(s1&0xFF) |
						rows[startR+(s2>>8)]>>(s2&0xFF) |
						rows[startR+(s3>>8)]>>(s3&0xFF) |
						rows[startR+(s4>>8)]>>(s4&0xFF) |
						rows[startR+(s5>>8)]>>(s5&0xFF)
				case 7:
					s0, s1, s2, s3, s4, s5, s6 := shifts[0], shifts[1], shifts[2], shifts[3], shifts[4], shifts[5], shifts[6]
					blocked = rows[startR+(s0>>8)]>>(s0&0xFF) |
						rows[startR+(s1>>8)]>>(s1&0xFF) |
						rows[startR+(s2>>8)]>>(s2&0xFF) |
						rows[startR+(s3>>8)]>>(s3&0xFF) |
						rows[startR+(s4>>8)]>>(s4&0xFF) |
						rows[startR+(s5>>8)]>>(s5&0xFF) |
						rows[startR+(s6>>8)]>>(s6&0xFF)
				default:
					for _, shift := range shifts {
						blocked |= rows[startR+(shift>>8)] >> (shift & 0xFF)
					}
				}

				// Find available positions
				available := (^blocked) & validMask
				if available == 0 {
					continue
				}

				// Try first available position - blocked mask should be exact
				startC := bits.TrailingZeros64(available)

				// Place directly without verification (blocked mask is exact)
				for i := 0; i < nRows; i++ {
					g.rows[startR+i] |= rowMasks[i] << startC
				}
				placed = true
				break orientLoop
			}
		}
		if !placed {
			return false
		}
	}
	return true
}

func (g *Grid) reset(w, h int) {
	g.width = w
	g.height = h
	if cap(g.rows) >= h {
		g.rows = g.rows[:h]
		for i := range g.rows {
			g.rows[i] = 0
		}
	} else {
		g.rows = make([]uint64, h)
	}
}

func canFit(allMasks [][]ShapeMask, region Region, cellCount int, shapes []ShapeEntry, g *Grid) bool {
	// Build shape list (reusing provided slice)
	shapes = shapes[:0]
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

	// Try greedy first (reusing grid)
	g.reset(region.width, region.height)
	if greedyPlace(g, allMasks, shapes) {
		return true
	}

	// Fall back to backtracking
	g.reset(region.width, region.height)
	solver := &Solver{
		grid:     g,
		allMasks: allMasks,
		shapes:   shapes,
	}

	skipsAllowed := gridArea - totalCells
	return solver.solve(0, skipsAllowed)
}

func part1Sequential(lines []string) int {
	allMasks, regions, cellCount := parseInput(lines)
	shapes := make([]ShapeEntry, 0, 300)
	g := &Grid{rows: make([]uint64, 64)}
	count := 0
	for _, r := range regions {
		if canFit(allMasks, r, cellCount, shapes, g) {
			count++
		}
	}
	return count
}

func part1(lines []string) int {
	allMasks, regions, cellCount := parseInput(lines)

	numWorkers := runtime.NumCPU()
	jobs := make(chan Region, len(regions))
	var wg sync.WaitGroup
	var count atomic.Int64

	// Start worker pool
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Pre-allocate for this worker
			shapes := make([]ShapeEntry, 0, 300)
			g := &Grid{rows: make([]uint64, 64)}
			for r := range jobs {
				if canFit(allMasks, r, cellCount, shapes, g) {
					count.Add(1)
				}
			}
		}()
	}

	// Send all jobs
	for _, region := range regions {
		jobs <- region
	}
	close(jobs)

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
