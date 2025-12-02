package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Direction int

const (
	Left Direction = iota
	Right
)

type Move struct {
	Direction Direction
	Amount    int
}

func parse(lines []string) []Move {
	moves := make([]Move, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var dir Direction
		switch line[0] {
		case 'L':
			dir = Left
		case 'R':
			dir = Right
		}
		amount, _ := strconv.Atoi(line[1:])
		moves = append(moves, Move{Direction: dir, Amount: amount})
	}
	return moves
}

func part1(moves []Move) int {
	dial := 50
	password := 0

	for _, m := range moves {
		switch m.Direction {
		case Left:
			dial -= m.Amount
		case Right:
			dial += m.Amount
		}

		// Wrap dial to 0-99 range
		dial = ((dial % 100) + 100) % 100

		if dial == 0 {
			password++
		}
	}

	return password
}

func part2(moves []Move) int {
	dial := 50
	password := 0

	for _, m := range moves {
		for range m.Amount {
			switch m.Direction {
			case Left:
				dial = (dial - 1 + 100) % 100
			case Right:
				dial = (dial + 1) % 100
			}

			if dial == 0 {
				password++
			}
		}
	}

	return password
}

func main() {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	moves := parse(lines)
	fmt.Println("Part 1:", part1(moves))
	fmt.Println("Part 2:", part2(moves))
}
