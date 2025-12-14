package main

import (
	"fmt"
	"testing"
)

// generateLargePolygon creates a rectangular polygon with approximately n red tiles
// It creates a rectangle with perimeter ~n, using the pattern:
// Start at (0,0), go right to (width,0), down to (width,height), left to (0,height), back to (0,0)
func generateLargePolygon(n int) []string {
	// For a rectangle with width w and height h, we have 4 corners
	// Perimeter tiles = 2*w + 2*h (corners counted in both)
	// To get ~n tiles, use width = n/4, height = n/4
	side := n / 4
	if side < 2 {
		side = 2
	}

	// Create rectangle corners: (0,0) -> (side,0) -> (side,side) -> (0,side)
	lines := []string{
		"0,0",
		fmt.Sprintf("%d,0", side),
		fmt.Sprintf("%d,%d", side, side),
		fmt.Sprintf("0,%d", side),
	}
	return lines
}

// generateComplexPolygon creates a more complex polygon with many corners
// This creates a zigzag pattern that results in more red tiles
func generateComplexPolygon(n int) []string {
	// Create a zigzag pattern going right then up/down alternating
	// Each "tooth" adds 4 corners
	teeth := n / 8
	if teeth < 10 {
		teeth = 10
	}

	lines := make([]string, 0, teeth*4)
	x, y := 0, 0

	for i := 0; i < teeth; i++ {
		// Bottom of tooth
		lines = append(lines, fmt.Sprintf("%d,%d", x, y))
		x += 2
		lines = append(lines, fmt.Sprintf("%d,%d", x, y))
		// Top of tooth
		y += 5
		lines = append(lines, fmt.Sprintf("%d,%d", x, y))
		x += 2
		lines = append(lines, fmt.Sprintf("%d,%d", x, y))
		y -= 5
	}

	// Close the polygon - go back along the top
	topY := 10
	lines = append(lines, fmt.Sprintf("%d,%d", x, topY))
	lines = append(lines, fmt.Sprintf("0,%d", topY))

	return lines
}

func TestPart1Example(t *testing.T) {
	input := []string{
		"7,1",
		"11,1",
		"11,7",
		"9,7",
		"9,5",
		"2,5",
		"2,3",
		"7,3",
	}
	got := part1(input)
	want := 50
	if got != want {
		t.Errorf("part1() = %d, want %d", got, want)
	}
}

func TestPart2Example(t *testing.T) {
	input := []string{
		"7,1",
		"11,1",
		"11,7",
		"9,7",
		"9,5",
		"2,5",
		"2,3",
		"7,3",
	}
	got := part2(input)
	want := 24
	if got != want {
		t.Errorf("part2() = %d, want %d", got, want)
	}
}

func BenchmarkPart2Small(b *testing.B) {
	// Simple rectangle - 4 corners
	input := generateLargePolygon(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		part2(input)
	}
}

func BenchmarkPart2Complex(b *testing.B) {
	// Complex zigzag polygon with many corners
	input := generateComplexPolygon(200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		part2(input)
	}
}

func BenchmarkPart2Large(b *testing.B) {
	// Larger complex polygon
	input := generateComplexPolygon(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		part2(input)
	}
}
