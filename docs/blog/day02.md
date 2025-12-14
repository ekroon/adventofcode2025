# Hunting for Repeating Patterns: A Tale of Digital Déjà Vu

When you're tasked with validating millions of ID numbers, pattern recognition becomes your best friend. Day 2 of Advent of Code 2025 challenged us to detect invalid IDs that exhibit suspicious repetition—numbers like `1212` (doubled sequences) or `123123123` (tripled patterns).

## The Pattern Detection Challenge

The puzzle splits into two parts:
- **Part 1**: Catch doubles—numbers where the digits repeat exactly twice (like `4747` or `123123`)
- **Part 2**: Catch all repeaters—any pattern that repeats 2 or more times (including `11`, `123123`, and `45454545`)

## Elegant Digit Extraction

The secret sauce lies in how we extract digits without string conversions. Using modulo and integer division, we peel digits off efficiently:

```go
func isInvalidDouble(id int) bool {
    var digits [20]byte
    n := 0
    for num := id; num > 0; num /= 10 {
        digits[n] = byte(num % 10)
        n++
    }
    // ... pattern matching logic
}
```

Each iteration, `num % 10` gives us the rightmost digit, then `num /= 10` shifts right. This reverse-order extraction is fast and allocation-free—critical when processing millions of numbers.

## Parallel Power with Goroutines

Processing large ID ranges sequentially would be painfully slow. Instead, we leverage Go's concurrency model:

1. Split each ID range into chunks based on `runtime.GOMAXPROCS(0)` (matching CPU cores)
2. Launch worker goroutines, each processing its chunk independently
3. Collect results through a buffered channel
4. Sum the partial results from all workers

```go
numWorkers := runtime.GOMAXPROCS(0)
results := make(chan int, numWorkers)
var wg sync.WaitGroup

// Launch workers for each chunk
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
```

This worker pool pattern distributes the computational load across all CPU cores, turning what could be a multi-second operation into a near-instantaneous scan. The beauty of Go's channels and goroutines makes concurrent programming feel natural and safe.

## Conclusion

By combining efficient digit manipulation with Go's powerful concurrency primitives, we built a solution that's both elegant and performant. Pattern detection meets parallel processing—a perfect harmony for tackling large-scale data validation challenges!
