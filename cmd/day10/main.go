package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	numLights int
	target    []bool  // target light pattern
	buttons   [][]int // each button is a list of light indices it toggles
}

func parseMachine(line string) Machine {
	// Format: [lights] (button1) (button2) ... {joltages}
	m := Machine{}

	// Extract [lights]
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	lights := line[start+1 : end]
	m.numLights = len(lights)
	m.target = make([]bool, m.numLights)
	for i, ch := range lights {
		m.target[i] = (ch == '#')
	}

	// Extract buttons (...)
	rest := line[end+1:]
	for {
		start := strings.Index(rest, "(")
		if start == -1 {
			break
		}
		end := strings.Index(rest, ")")
		buttonStr := rest[start+1 : end]

		// Parse comma-separated indices
		button := []int{}
		for _, s := range strings.Split(buttonStr, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				idx, _ := strconv.Atoi(s)
				button = append(button, idx)
			}
		}
		m.buttons = append(m.buttons, button)
		rest = rest[end+1:]
	}

	return m
}

// Gaussian elimination over GF(2) to solve the system, finding minimum 1s in solution
func solveGF2(m Machine) int {
	numButtons := len(m.buttons)
	numLights := m.numLights

	// Create augmented matrix [A|b] where A[i][j] = 1 if button j toggles light i
	// b[i] = target state of light i
	matrix := make([][]int, numLights)
	for i := range matrix {
		matrix[i] = make([]int, numButtons+1)
		if m.target[i] {
			matrix[i][numButtons] = 1
		}
	}

	// Fill in the button effects: matrix[light][button] = 1 if button toggles light
	for btnIdx, button := range m.buttons {
		for _, lightIdx := range button {
			matrix[lightIdx][btnIdx] = 1
		}
	}

	// Gaussian elimination to reduced row echelon form
	pivotRow := 0
	pivotColForRow := make([]int, numLights) // pivotColForRow[row] = column of pivot, or -1
	for i := range pivotColForRow {
		pivotColForRow[i] = -1
	}
	isPivotCol := make([]bool, numButtons)

	for col := 0; col < numButtons && pivotRow < numLights; col++ {
		// Find pivot
		foundRow := -1
		for row := pivotRow; row < numLights; row++ {
			if matrix[row][col] == 1 {
				foundRow = row
				break
			}
		}

		if foundRow == -1 {
			continue
		}

		// Swap rows
		matrix[pivotRow], matrix[foundRow] = matrix[foundRow], matrix[pivotRow]
		pivotColForRow[pivotRow] = col
		isPivotCol[col] = true

		// Eliminate
		for row := 0; row < numLights; row++ {
			if row != pivotRow && matrix[row][col] == 1 {
				for c := 0; c <= numButtons; c++ {
					matrix[row][c] ^= matrix[pivotRow][c]
				}
			}
		}
		pivotRow++
	}

	// Check for inconsistency
	for row := 0; row < numLights; row++ {
		allZero := true
		for col := 0; col < numButtons; col++ {
			if matrix[row][col] == 1 {
				allZero = false
				break
			}
		}
		if allZero && matrix[row][numButtons] == 1 {
			return -1 // No solution
		}
	}

	// Identify free variables (columns without pivots)
	freeVars := []int{}
	for col := 0; col < numButtons; col++ {
		if !isPivotCol[col] {
			freeVars = append(freeVars, col)
		}
	}

	// Try all 2^k combinations of free variables to find minimum solution
	numFree := len(freeVars)
	minCount := -1

	for mask := 0; mask < (1 << numFree); mask++ {
		// Set free variables according to mask
		solution := make([]int, numButtons)
		for i, col := range freeVars {
			if (mask & (1 << i)) != 0 {
				solution[col] = 1
			}
		}

		// Back-substitute to find pivot variable values
		for row := numLights - 1; row >= 0; row-- {
			if pivotColForRow[row] == -1 {
				continue
			}
			col := pivotColForRow[row]
			val := matrix[row][numButtons]
			// Subtract contributions from variables to the right
			for c := col + 1; c < numButtons; c++ {
				val ^= matrix[row][c] * solution[c]
			}
			solution[col] = val
		}

		// Count number of 1s in solution
		count := 0
		for _, v := range solution {
			count += v
		}

		if minCount == -1 || count < minCount {
			minCount = count
		}
	}

	return minCount
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		machine := parseMachine(line)
		presses := solveGF2(machine)
		if presses >= 0 {
			total += presses
		}
	}
	return total
}

func part2(lines []string) int {
	// TODO: implement
	return 0
}

func main() {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println("Part 1:", part1(lines))
	fmt.Println("Part 2:", part2(lines))
}
