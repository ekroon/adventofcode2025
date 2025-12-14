package main

import (
	"bufio"
	"os"
	"testing"
)

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

var exampleInput2 = []string{
	"svr: aaa bbb",
	"aaa: fft",
	"fft: ccc",
	"bbb: tty",
	"tty: ccc",
	"ccc: ddd eee",
	"ddd: hub",
	"hub: fff",
	"eee: dac",
	"dac: fff",
	"fff: ggg hhh",
	"ggg: out",
	"hhh: out",
}

func TestPart2(t *testing.T) {
	got := part2(exampleInput2)
	want := 2
	if got != want {
		t.Errorf("part2() = %d, want %d", got, want)
	}
}

func BenchmarkPart1(b *testing.B) {
	for b.Loop() {
		part1(exampleInput)
	}
}

func BenchmarkPart2(b *testing.B) {
	for b.Loop() {
		part2(exampleInput2)
	}
}

func loadInput(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func BenchmarkPart1Real(b *testing.B) {
	lines := loadInput("../../inputs/day11.txt")
	b.ResetTimer()
	for b.Loop() {
		part1(lines)
	}
}

func BenchmarkPart2Real(b *testing.B) {
	lines := loadInput("../../inputs/day11.txt")
	b.ResetTimer()
	for b.Loop() {
		part2(lines)
	}
}
