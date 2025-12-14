package main

import "testing"

var testInput = []string{
	".......S.......",
	"...............",
	".......^.......",
	"...............",
	"......^.^......",
	"...............",
	".....^.^.^.....",
	"...............",
	"....^.^...^....",
	"...............",
	"...^.^...^.^...",
	"...............",
	"..^...^.....^..",
	"...............",
	".^.^.^.^.^...^.",
	"...............",
}

func TestPart1(t *testing.T) {
	got := part1(testInput)
	want := 21
	if got != want {
		t.Errorf("part1() = %d, want %d", got, want)
	}
}

func TestPart2(t *testing.T) {
	got := part2(testInput)
	want := 40
	if got != want {
		t.Errorf("part2() = %d, want %d", got, want)
	}
}
