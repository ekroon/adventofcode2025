package main

import (
	"testing"
)

func TestParseMachine(t *testing.T) {
	line := "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}"
	m := parseMachine(line)

	t.Logf("Machine: numLights=%d, target=%v, buttons=%v", m.numLights, m.target, m.buttons)

	if m.numLights != 4 {
		t.Errorf("Expected 4 lights, got %d", m.numLights)
	}

	// [.##.] means lights 1 and 2 should be on (true), 0 and 3 off (false)
	expectedTarget := []bool{false, true, true, false}
	for i, v := range expectedTarget {
		if m.target[i] != v {
			t.Errorf("target[%d]: expected %v, got %v", i, v, m.target[i])
		}
	}
}

func TestSolveGF2(t *testing.T) {
	tests := []struct {
		line     string
		expected int
	}{
		{"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}", -1}, // placeholder, need to verify
		{"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}", -1},
		{"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}", -1},
	}

	for _, tt := range tests {
		m := parseMachine(tt.line)
		t.Logf("Machine: numLights=%d, target=%v, buttons=%v", m.numLights, m.target, m.buttons)
		result := solveGF2(m)
		t.Logf("Line: %s => result: %d", tt.line, result)
	}
}

func TestPart1(t *testing.T) {
	input := []string{
		"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
		"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
		"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
	}

	result := part1(input)
	t.Logf("Part 1 result: %d (expected 7)", result)

	if result != 7 {
		t.Errorf("Expected 7, got %d", result)
	}
}
