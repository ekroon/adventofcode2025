package main

import "testing"

func TestPart2(t *testing.T) {
	input := []string{
		"123 328  51 64 ",
		" 45 64  387 23 ",
		"  6 98  215 314",
		"*   +   *   +  ",
	}

	// Debug: print what we're reading
	t.Logf("Line lengths: %d, %d, %d, %d", len(input[0]), len(input[1]), len(input[2]), len(input[3]))
	t.Logf("Line 3: %q", input[3])

	want := 3263827
	got := part2(input)

	if got != want {
		t.Errorf("part2() = %d, want %d", got, want)
	}
}

func TestPart2Debug(t *testing.T) {
	input := []string{
		"123 328  51 64 ",
		" 45 64  387 23 ",
		"  6 98  215 314",
		"*   +   *   +  ",
	}

	// Trace through manually
	maxWidth := 15
	for col := maxWidth - 1; col >= 0; col-- {
		for row := 0; row < 4; row++ {
			line := input[row]
			if col < len(line) {
				t.Logf("col=%d row=%d char=%q", col, row, line[col])
			} else {
				t.Logf("col=%d row=%d OUT OF BOUNDS (len=%d)", col, row, len(line))
			}
		}
	}
}
