package main

import (
	"os"
	"strings"
	"testing"
)

func loadInput(b *testing.B) []string {
	b.Helper()
	data, err := os.ReadFile("../../inputs/day05.txt")
	if err != nil {
		b.Skip("inputs/day05.txt not found")
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

func BenchmarkPart1(b *testing.B) {
	lines := loadInput(b)
	b.ResetTimer()
	for b.Loop() {
		part1(lines)
	}
}

func BenchmarkPart2(b *testing.B) {
	lines := loadInput(b)
	b.ResetTimer()
	for b.Loop() {
		part2(lines)
	}
}
