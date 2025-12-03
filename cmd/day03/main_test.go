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
