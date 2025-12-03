package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// isInvalidDouble checks if a number consists of a sequence repeated exactly twice
func isInvalidDouble(id int) bool {
	var digits [20]byte
	n := 0
	for num := id; num > 0; num /= 10 {
		digits[n] = byte(num % 10)
		n++
	}
	if n%2 != 0 {
		return false
	}
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	half := n / 2
	for i := 0; i < half; i++ {
		if digits[i] != digits[i+half] {
			return false
		}
	}
	return true
}

// isInvalidRepeated checks if a number consists of a sequence repeated 2+ times
func isInvalidRepeated(id int) bool {
	var digits [20]byte
	n := 0
	for num := id; num > 0; num /= 10 {
		digits[n] = byte(num % 10)
		n++
	}
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	for patLen := 1; patLen <= n/2; patLen++ {
		if n%patLen != 0 {
			continue
		}
		match := true
		for i := patLen; i < n; i++ {
			if digits[i] != digits[i%patLen] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func part1(lines []string) int {
	if len(lines) == 0 {
		return 0
	}

	type idRange struct{ start, end int }
	var ranges []idRange

	for r := range strings.SplitSeq(lines[0], ",") {
		startStr, endStr, _ := strings.Cut(r, "-")
		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)
		ranges = append(ranges, idRange{start, end})
	}

	numWorkers := runtime.GOMAXPROCS(0)
	results := make(chan int, numWorkers)
	var wg sync.WaitGroup

	for _, rng := range ranges {
		rangeSize := rng.end - rng.start + 1
		chunkSize := (rangeSize + numWorkers - 1) / numWorkers

		for w := 0; w < numWorkers; w++ {
			chunkStart := rng.start + w*chunkSize
			chunkEnd := min(chunkStart+chunkSize-1, rng.end)
			if chunkStart > rng.end {
				break
			}

			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				localSum := 0
				for id := start; id <= end; id++ {
					if isInvalidDouble(id) {
						localSum += id
					}
				}
				results <- localSum
			}(chunkStart, chunkEnd)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for partial := range results {
		sum += partial
	}
	return sum
}

func part2(lines []string) int {
	if len(lines) == 0 {
		return 0
	}

	type idRange struct{ start, end int }
	var ranges []idRange

	for r := range strings.SplitSeq(lines[0], ",") {
		startStr, endStr, _ := strings.Cut(r, "-")
		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)
		ranges = append(ranges, idRange{start, end})
	}

	numWorkers := runtime.GOMAXPROCS(0)
	results := make(chan int, numWorkers)
	var wg sync.WaitGroup

	for _, rng := range ranges {
		rangeSize := rng.end - rng.start + 1
		chunkSize := (rangeSize + numWorkers - 1) / numWorkers

		for w := 0; w < numWorkers; w++ {
			chunkStart := rng.start + w*chunkSize
			chunkEnd := min(chunkStart+chunkSize-1, rng.end)
			if chunkStart > rng.end {
				break
			}

			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				localSum := 0
				for id := start; id <= end; id++ {
					if isInvalidRepeated(id) {
						localSum += id
					}
				}
				results <- localSum
			}(chunkStart, chunkEnd)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for partial := range results {
		sum += partial
	}
	return sum
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
