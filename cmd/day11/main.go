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

func countPaths(graph map[string][]string, current, target string) int {
	if current == target {
		return 1
	}
	neighbors, ok := graph[current]
	if !ok {
		return 0
	}
	count := 0
	for _, next := range neighbors {
		count += countPaths(graph, next, target)
	}
	return count
}

func part1(lines []string) int {
	graph := parseGraph(lines)
	return countPaths(graph, "you", "out")
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
