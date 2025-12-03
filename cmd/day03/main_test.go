package main

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func generateLines(numLines, lineLength int) []string {
	lines := make([]string, numLines)
	var sb strings.Builder
	for i := range lines {
		sb.Reset()
		for j := 0; j < lineLength; j++ {
			sb.WriteByte(byte('0' + rand.Intn(10)))
		}
		lines[i] = sb.String()
	}
	return lines
}

func BenchmarkPart1(b *testing.B) {
	lines := generateLines(200, 100)
	b.ResetTimer()
	for b.Loop() {
		part1(lines)
	}
}

func BenchmarkPart2(b *testing.B) {
	lines := generateLines(200, 100)
	b.ResetTimer()
	for b.Loop() {
		part2(lines)
	}
}

func TestMaxNumber(t *testing.T) {
	tests := []struct {
		line   string
		count  int
		expect int
	}{
		{"1234", 2, 34},
		{"9123311181", 2, 98},
		{"2213322222223222132231233322432224423222226232323522232215252332221122222231232224232722131522422232", 12, 731522422232},
	}
	for _, tc := range tests {
		t.Run(tc.line[:min(10, len(tc.line))]+"_"+strconv.Itoa(tc.count), func(t *testing.T) {
			got := maxNumber(tc.line, tc.count)
			if got != tc.expect {
				t.Errorf("maxNumber(%q, %d) = %d, want %d", tc.line, tc.count, got, tc.expect)
			}
		})
	}
}

// maxNumberNaive is the brute-force approach: try all combinations of `count` indices
func maxNumberNaive(line string, count int) int {
	maxNum := 0
	var search func(start, depth, current int)
	search = func(start, depth, current int) {
		if depth == count {
			maxNum = max(maxNum, current)
			return
		}
		for i := start; i <= len(line)-(count-depth); i++ {
			search(i+1, depth+1, current*10+int(line[i]-'0'))
		}
	}
	search(0, 0, 0)
	return maxNum
}

func BenchmarkPart2Naive(b *testing.B) {
	lines := generateLines(1, 30) // Short line - naive is extremely slow
	b.ResetTimer()
	for b.Loop() {
		total := 0
		for _, line := range lines {
			total += maxNumberNaive(line, 12)
		}
	}
}

func BenchmarkPart2Greedy(b *testing.B) {
	lines := generateLines(1, 30) // Same short line for fair comparison
	b.ResetTimer()
	for b.Loop() {
		total := 0
		for _, line := range lines {
			total += maxNumber(line, 12)
		}
	}
}
