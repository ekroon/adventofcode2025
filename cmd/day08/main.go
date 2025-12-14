package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Edge struct {
	i, j   int
	distSq int
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

// Union-Find with path compression and union by rank
type UnionFind struct {
	parent []int
	rank   []int
	size   []int
}

func newUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	size := make([]int, n)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent, rank, size}
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
	if uf.rank[px] < uf.rank[py] {
		px, py = py, px
	}
	uf.parent[py] = px
	uf.size[px] += uf.size[py]
	if uf.rank[px] == uf.rank[py] {
		uf.rank[px]++
	}
	return true
}

func part1(lines []string) int {
	points := parsePoints(lines)
	n := len(points)

	// Generate all edges with distances
	edges := make([]Edge, 0, n*(n-1)/2)
	for i := range n {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{i, j, distanceSquared(points[i], points[j])})
		}
	}

	// Sort by distance
	slices.SortFunc(edges, func(a, b Edge) int {
		return cmp.Compare(a.distSq, b.distSq)
	})

	// Connect 1000 closest pairs
	uf := newUnionFind(n)
	connections := 0
	for _, e := range edges {
		if connections >= 1000 {
			break
		}
		uf.union(e.i, e.j)
		connections++
	}

	// Find circuit sizes
	sizes := make(map[int]int)
	for i := range n {
		root := uf.find(i)
		sizes[root] = uf.size[root]
	}

	// Get unique sizes and sort descending
	sizeList := make([]int, 0, len(sizes))
	for _, s := range sizes {
		sizeList = append(sizeList, s)
	}
	slices.SortFunc(sizeList, func(a, b int) int {
		return cmp.Compare(b, a) // descending
	})

	// Multiply top 3
	return sizeList[0] * sizeList[1] * sizeList[2]
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
