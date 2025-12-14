package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Edge struct {
	distSq int
	i, j   uint16
}

func parsePoints(lines []string) []Point {
	points := make([]Point, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		points = append(points, Point{x, y, z})
	}
	return points
}

func distanceSquared(a, b Point) int {
	dx, dy, dz := a.x-b.x, a.y-b.y, a.z-b.z
	return dx*dx + dy*dy + dz*dz
}

// EdgeHeap implements a min-heap of edges by distance
type EdgeHeap []Edge

func (h EdgeHeap) Len() int           { return len(h) }
func (h EdgeHeap) Less(i, j int) bool { return h[i].distSq < h[j].distSq }
func (h EdgeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *EdgeHeap) Push(x any)        { *h = append(*h, x.(Edge)) }
func (h *EdgeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// pop returns and removes the minimum edge without interface allocation
func (h *EdgeHeap) pop() Edge {
	e := (*h)[0]
	n := len(*h) - 1
	(*h)[0] = (*h)[n]
	*h = (*h)[:n]
	if n > 0 {
		h.down(0)
	}
	return e
}

func (h EdgeHeap) down(i int) {
	for {
		left := 2*i + 1
		if left >= len(h) {
			break
		}
		j := left
		if right := left + 1; right < len(h) && h[right].distSq < h[left].distSq {
			j = right
		}
		if h[i].distSq <= h[j].distSq {
			break
		}
		h[i], h[j] = h[j], h[i]
		i = j
	}
}

// init establishes heap ordering
func (h EdgeHeap) init() {
	for i := len(h)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

// buildEdgeHeap parses points and returns them with all edges in a min-heap
func buildEdgeHeap(lines []string) ([]Point, *EdgeHeap) {
	points := parsePoints(lines)
	n := len(points)

	edges := make(EdgeHeap, 0, n*(n-1)/2)
	for i := range n {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{distanceSquared(points[i], points[j]), uint16(i), uint16(j)})
		}
	}
	edges.init()

	return points, &edges
}

// Union-Find with path compression and union by size
type UnionFind struct {
	parent []int
	size   []int
}

func newUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	size := make([]int, n)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent, size}
}

func (uf *UnionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) union(x, y int) bool {
	px, py := uf.find(x), uf.find(y)
	if px == py {
		return false // already in same circuit
	}
	// Union by size: attach smaller tree under larger
	if uf.size[px] < uf.size[py] {
		px, py = py, px
	}
	uf.parent[py] = px
	uf.size[px] += uf.size[py]
	return true
}

func part1(lines []string) int {
	points, edges := buildEdgeHeap(lines)
	n := len(points)

	// Connect 1000 closest pairs
	uf := newUnionFind(n)
	connections := 0
	for edges.Len() > 0 && connections < 1000 {
		e := edges.pop()
		uf.union(int(e.i), int(e.j))
		connections++
	}

	// Collect circuit sizes
	sizes := make(map[int]int)
	for i := range n {
		root := uf.find(i)
		sizes[root] = uf.size[root]
	}

	// Sort descending and multiply top 3
	sizeList := slices.Collect(maps.Values(sizes))
	slices.Sort(sizeList)
	slices.Reverse(sizeList)

	return sizeList[0] * sizeList[1] * sizeList[2]
}

func part2(lines []string) int {
	points, edges := buildEdgeHeap(lines)
	n := len(points)

	// Connect until all in one circuit
	uf := newUnionFind(n)
	var lastEdge Edge
	for edges.Len() > 0 {
		e := edges.pop()
		if uf.union(int(e.i), int(e.j)) {
			lastEdge = e
			// Check if all connected (root's size equals n)
			if uf.size[uf.find(int(e.i))] == n {
				break
			}
		}
	}

	// Multiply X coordinates of last connected pair
	return points[lastEdge.i].x * points[lastEdge.j].x
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
