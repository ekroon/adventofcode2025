package main

import "testing"

var exampleInput = []string{
	"aaa: you hhh",
	"you: bbb ccc",
	"bbb: ddd eee",
	"ccc: ddd eee fff",
	"ddd: ggg",
	"eee: out",
	"fff: out",
	"ggg: out",
	"hhh: ccc fff iii",
	"iii: out",
}

func TestPart1(t *testing.T) {
	got := part1(exampleInput)
	want := 5
	if got != want {
		t.Errorf("part1() = %d, want %d", got, want)
	}
}
