# Day 05: The Art of Merging Intervals

When life gives you overlapping ranges, you merge them! Day 05 of Advent of Code presented a classic computational geometry problem: given a collection of number ranges like "10-20" and "15-30", how do we efficiently determine if specific numbers fall within the combined coverage?

## The Range Merging Dance

The heart of this solution lies in a beautiful interval merging algorithm. The idea is elegantly simple: sort all ranges by their starting point, then walk through them, merging any that overlap or sit right next to each other.

Here's the magic in action:

```go
func mergeRanges(ranges []Range) []Range {
    // Sort by low bound
    slices.SortFunc(sorted, func(a, b Range) int {
        return a.low - b.low
    })

    merged := []Range{sorted[0]}
    for _, r := range sorted[1:] {
        last := &merged[len(merged)-1]
        if r.low <= last.high+1 {
            // Overlapping or adjacent, extend the range
            last.high = max(last.high, r.high)
        } else {
            // Gap, start new range
            merged = append(merged, r)
        }
    }
    return merged
}
```

## The Adjacent Range Gotcha

Notice that critical `+1` in the merge condition? This handles adjacent ranges that are only one apart, like [1-5] and [6-10]. These should merge into [1-10] since there's no actual gap in coverage. Missing this detail would leave you with fragmented ranges when you actually have continuous coverage!

## Binary Search: The Speed Boost

Once we have merged ranges, Part 1 asks us to check hundreds of numbers against these ranges. A naive approach would scan every range for each number—yawn. Instead, we leverage binary search since our merged ranges are already sorted:

```go
func inMergedRanges(n int, merged []Range) bool {
    // Binary search: find the first range where low > n
    i, _ := slices.BinarySearchFunc(merged, n, func(r Range, target int) int {
        return r.low - target
    })
    // Check the range before (if exists) - it could contain n
    if i > 0 && n <= merged[i-1].high {
        return true
    }
    // Check if we landed exactly on a range that starts with n
    if i < len(merged) && n >= merged[i].low && n <= merged[i].high {
        return true
    }
    return false
}
```

This drops our lookup from O(n) to O(log n)—a massive win when checking many numbers!

Part 2 was almost a victory lap: once you have merged ranges, just sum up their sizes. The hard work was already done in the merging algorithm.

**Key Takeaway**: When dealing with intervals, always consider merging first. It simplifies queries and often reveals elegant solutions hiding in plain sight.
