package main

import (
	"os"
	"strings"
	"testing"
)

var example = `0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2`

func TestPart1(t *testing.T) {
	lines := strings.Split(example, "\n")
	// Use backtracking for small examples since arithmetic is an approximation
	got := part1Sequential(lines)
	want := 2
	if got != want {
		t.Errorf("part1() = %d, want %d", got, want)
	}
}

func loadInput() []string {
	data, err := os.ReadFile("../../inputs/day12.txt")
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

func BenchmarkPart1(b *testing.B) {
	lines := loadInput()
	if lines == nil {
		b.Skip("input file not found")
	}
	b.ResetTimer()
	for range b.N {
		part1(lines)
	}
}

func BenchmarkPart1Sequential(b *testing.B) {
	lines := loadInput()
	if lines == nil {
		b.Skip("input file not found")
	}
	b.ResetTimer()
	for range b.N {
		part1Sequential(lines)
	}
}

func BenchmarkPart1Arithmetic(b *testing.B) {
	lines := loadInput()
	if lines == nil {
		b.Skip("input file not found")
	}
	b.ResetTimer()
	for range b.N {
		part1Arithmetic(lines)
	}
}

func TestArithmeticVsActual(t *testing.T) {
	lines := loadInput()
	if lines == nil {
		t.Skip("input file not found")
	}
	actual := part1Sequential(lines)
	arithmetic := part1Arithmetic(lines)
	t.Logf("Actual: %d, Arithmetic: %d, Diff: %d", actual, arithmetic, arithmetic-actual)
}
