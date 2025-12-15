# Greedy for Digits: The Art of Maximum Number Selection

When faced with a string of digits and asked to pick just a few to form the largest possible number, your first instinct might be to try every combination. But what if I told you there's a clever greedy approach that's not just faster, but actually guarantees the optimal answer?

## The Challenge

Day 03's puzzle presents a deceptively simple problem: given a string of digits, select N digits **in order** to form the maximum possible number. Part 1 asks for 2 digits, while Part 2 cranks it up to 12. The constraint that digits must maintain their relative order is crucial—it transforms this from a simple "pick the largest digits" problem into something more nuanced.

## Why Greedy Works

The beauty of this problem lies in its **optimal substructure**: if you've picked the best first digit, the remaining problem is just picking the best digits from what's left. This is the hallmark of problems where greedy algorithms shine.

The key insight? For each position in your result, pick the **largest available digit** that still leaves enough digits behind it for the remaining positions. If you need to pick 12 digits total and you're selecting the 3rd one, you must leave at least 9 more digits after your choice.

## The Clever Bounds Calculation

Here's where the magic happens in the code:

```go
for i := 0; i < count; i++ {
    // We need (count - i - 1) more digits after this one
    // So we can pick from start to n-(count-i-1)-1 = n-count+i
    end := n - count + i
    bestDigit := byte('0')
    bestIdx := start
    for j := start; j <= end; j++ {
        if line[j] > bestDigit {
            bestDigit = line[j]
            bestIdx = j
        }
    }
    result = result*10 + int(bestDigit-'0')
    start = bestIdx + 1
}
```

The formula `end = n - count + i` elegantly ensures we never run out of digits. For position `i`, we need `count - i - 1` more digits after our selection, so the latest we can pick is at index `n - (count - i - 1) - 1`, which simplifies to `n - count + i`. Beautiful!

This greedy approach transforms what could be trillions of combinations into a simple O(n × k) solution—taking a 100-character line from computationally infeasible to solved in microseconds.
