package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseGraph(lines []string) map[string][]string {
	graph := make(map[string][]string)
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}
		from := parts[0]
		targets := strings.Fields(parts[1])
		graph[from] = targets
	}
	return graph
}

func countPaths(graph map[string][]string, current, target string, cache map[string]int) int {
	if current == target {
		return 1
	}
	if v, ok := cache[current]; ok {
		return v
	}
	neighbors, ok := graph[current]
	if !ok {
		return 0
	}
	count := 0
	for _, next := range neighbors {
		count += countPaths(graph, next, target, cache)
	}
	cache[current] = count
	return count
}

func part1(lines []string) int {
	graph := parseGraph(lines)
	cache := make(map[string]int)
	return countPaths(graph, "you", "out", cache)
}

type cacheKey struct {
	node       string
	visitedDac bool
	visitedFft bool
}

func countPathsWithRequired(graph map[string][]string, current, target string, visitedDac, visitedFft bool, cache map[cacheKey]int) int {
	if current == "dac" {
		visitedDac = true
	}
	if current == "fft" {
		visitedFft = true
	}
	key := cacheKey{current, visitedDac, visitedFft}
	if v, ok := cache[key]; ok {
		return v
	}
	if current == target {
		if visitedDac && visitedFft {
			return 1
		}
		return 0
	}
	neighbors, ok := graph[current]
	if !ok {
		return 0
	}
	count := 0
	for _, next := range neighbors {
		count += countPathsWithRequired(graph, next, target, visitedDac, visitedFft, cache)
	}
	cache[key] = count
	return count
}

func part2(lines []string) int {
	graph := parseGraph(lines)
	cache := make(map[cacheKey]int)
	return countPathsWithRequired(graph, "svr", "out", false, false, cache)
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
