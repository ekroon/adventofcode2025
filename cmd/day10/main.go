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
	joltages  []int   // target joltage levels for part2
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

	// Extract joltages {...}
	jStart := strings.Index(line, "{")
	jEnd := strings.Index(line, "}")
	if jStart != -1 && jEnd != -1 {
		joltageStr := line[jStart+1 : jEnd]
		for _, s := range strings.Split(joltageStr, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				val, _ := strconv.Atoi(s)
				m.joltages = append(m.joltages, val)
			}
		}
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
	total := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		machine := parseMachine(line)
		presses := solveJoltage(machine)
		if presses >= 0 {
			total += presses
		}
	}
	return total
}

// solveJoltage finds minimum button presses to achieve target joltages
// This is an integer linear programming problem:
// Minimize sum(x_i) subject to A*x = b where x_i >= 0
func solveJoltage(m Machine) int {
	numButtons := len(m.buttons)
	numCounters := len(m.joltages)

	// Build matrix A where A[counter][button] = 1 if button affects counter
	// We need to solve A*x = joltages with x >= 0, minimizing sum(x)
	affects := make([][]bool, numCounters)
	for i := range affects {
		affects[i] = make([]bool, numButtons)
	}
	for btnIdx, button := range m.buttons {
		for _, counterIdx := range button {
			if counterIdx < numCounters {
				affects[counterIdx][btnIdx] = true
			}
		}
	}

	// Use Gaussian elimination over rationals to find solution space
	// Then search for minimum non-negative integer solution

	// Create augmented matrix [A|b] with rational arithmetic
	type frac struct{ n, d int } // numerator/denominator
	gcd := func(a, b int) int {
		if a < 0 {
			a = -a
		}
		if b < 0 {
			b = -b
		}
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	simplify := func(f frac) frac {
		if f.n == 0 {
			return frac{0, 1}
		}
		g := gcd(f.n, f.d)
		f.n /= g
		f.d /= g
		if f.d < 0 {
			f.n, f.d = -f.n, -f.d
		}
		return f
	}
	sub := func(a, b frac) frac {
		return simplify(frac{a.n*b.d - b.n*a.d, a.d * b.d})
	}
	mul := func(a, b frac) frac {
		return simplify(frac{a.n * b.n, a.d * b.d})
	}
	div := func(a, b frac) frac {
		return simplify(frac{a.n * b.d, a.d * b.n})
	}

	matrix := make([][]frac, numCounters)
	for i := range matrix {
		matrix[i] = make([]frac, numButtons+1)
		for j := 0; j < numButtons; j++ {
			if affects[i][j] {
				matrix[i][j] = frac{1, 1}
			} else {
				matrix[i][j] = frac{0, 1}
			}
		}
		matrix[i][numButtons] = frac{m.joltages[i], 1}
	}

	// Gaussian elimination to RREF
	pivotRow := 0
	pivotColForRow := make([]int, numCounters)
	for i := range pivotColForRow {
		pivotColForRow[i] = -1
	}
	isPivotCol := make([]bool, numButtons)

	for col := 0; col < numButtons && pivotRow < numCounters; col++ {
		// Find non-zero pivot
		foundRow := -1
		for row := pivotRow; row < numCounters; row++ {
			if matrix[row][col].n != 0 {
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

		// Scale pivot row to make pivot = 1
		scale := matrix[pivotRow][col]
		for c := 0; c <= numButtons; c++ {
			matrix[pivotRow][c] = div(matrix[pivotRow][c], scale)
		}

		// Eliminate column in other rows
		for row := 0; row < numCounters; row++ {
			if row != pivotRow && matrix[row][col].n != 0 {
				factor := matrix[row][col]
				for c := 0; c <= numButtons; c++ {
					matrix[row][c] = sub(matrix[row][c], mul(factor, matrix[pivotRow][c]))
				}
			}
		}
		pivotRow++
	}

	// Check for inconsistency (row with all zeros on left but non-zero on right)
	for row := 0; row < numCounters; row++ {
		allZero := true
		for col := 0; col < numButtons; col++ {
			if matrix[row][col].n != 0 {
				allZero = false
				break
			}
		}
		if allZero && matrix[row][numButtons].n != 0 {
			return -1 // No solution
		}
	}

	// Identify free variables
	freeVars := []int{}
	for col := 0; col < numButtons; col++ {
		if !isPivotCol[col] {
			freeVars = append(freeVars, col)
		}
	}

	// Find minimum non-negative integer solution by searching over free variables
	// For each free variable, we need to find bounds that keep all variables >= 0
	// This is a bounded search - we iterate through reasonable values

	numFree := len(freeVars)
	minTotal := -1

	// Determine reasonable upper bound for free variables
	maxVal := 0
	for _, j := range m.joltages {
		if j > maxVal {
			maxVal = j
		}
	}

	// Recursive search over free variable values
	var search func(idx int, freeValues []int)
	search = func(idx int, freeValues []int) {
		if idx == numFree {
			// Compute pivot variable values
			solution := make([]frac, numButtons)
			for i := range solution {
				solution[i] = frac{0, 1}
			}
			for i, col := range freeVars {
				solution[col] = frac{freeValues[i], 1}
			}

			// Back-substitute
			valid := true
			for row := numCounters - 1; row >= 0; row-- {
				if pivotColForRow[row] == -1 {
					continue
				}
				col := pivotColForRow[row]
				val := matrix[row][numButtons]
				for c := col + 1; c < numButtons; c++ {
					val = sub(val, mul(matrix[row][c], solution[c]))
				}
				solution[col] = val
				// Check if integer and non-negative
				if val.d != 1 || val.n < 0 {
					valid = false
					break
				}
			}

			if valid {
				// Check all values are non-negative integers
				total := 0
				for _, v := range solution {
					if v.d != 1 || v.n < 0 {
						valid = false
						break
					}
					total += v.n
				}
				if valid && (minTotal == -1 || total < minTotal) {
					minTotal = total
				}
			}
			return
		}

		// Try values for this free variable
		for v := 0; v <= maxVal; v++ {
			freeValues[idx] = v
			search(idx+1, freeValues)
		}
	}

	search(0, make([]int, numFree))
	return minTotal
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
