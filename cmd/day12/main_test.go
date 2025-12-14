package main

import (
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
	got := part1(lines)
	want := 2
	if got != want {
		t.Errorf("part1() = %d, want %d", got, want)
	}
}
