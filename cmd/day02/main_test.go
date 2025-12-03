package main

import (
	"os"
	"strings"
	"testing"
)

func BenchmarkIsInvalidDouble(b *testing.B) {
	testCases := []int{
		11,       // 2 digits, invalid
		123,      // 3 digits, valid
		1234,     // 4 digits, valid
		123123,   // 6 digits, invalid
		1234567,  // 7 digits, valid
		12341234, // 8 digits, invalid
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			isInvalidDouble(tc)
		}
	}
}

func BenchmarkIsInvalidRepeated(b *testing.B) {
	testCases := []int{
		11,       // 2 digits, invalid
		123,      // 3 digits, valid
		1234,     // 4 digits, valid
		123123,   // 6 digits, invalid
		1234567,  // 7 digits, valid
		12341234, // 8 digits, invalid
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			isInvalidRepeated(tc)
		}
	}
}

func loadInput(b *testing.B) []string {
	b.Helper()
	data, err := os.ReadFile("../../inputs/day02.txt")
	if err != nil {
		b.Skip("inputs/day02.txt not found")
	}
	return []string{strings.TrimSpace(string(data))}
}

func BenchmarkPart1RealInput(b *testing.B) {
	lines := loadInput(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		part1(lines)
	}
}

func BenchmarkPart2RealInput(b *testing.B) {
	lines := loadInput(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		part2(lines)
	}
}
