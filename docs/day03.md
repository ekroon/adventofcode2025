# Day 03: Battery Banks

## Problem

Given lines of single digits, find the maximum number that can be formed by picking digits in order.

- **Part 1**: Pick 2 digits
- **Part 2**: Pick 12 digits

## Algorithm: Greedy Selection

The naive approach tries all combinations - for picking k digits from n characters, this is C(n,k) combinations.

For Part 2 with 100-character lines and k=12:
- C(100,12) ≈ **1.05 trillion** combinations per line

Instead, we use a **greedy algorithm**:

1. For each position i (0 to k-1), pick the largest digit that still leaves enough remaining digits
2. The search window is `[start, n-k+i]` where `start` is the index after the previous pick

### Example

Line: `9123311181`, pick 2 digits:

```
Step 1: Search [0,8], find '9' at index 0
Step 2: Search [1,9], find '8' at index 9
Result: 98
```

Line: `1234`, pick 2 digits:

```
Step 1: Search [0,2], find '3' at index 2
Step 2: Search [3,3], find '4' at index 3
Result: 34
```

## Complexity

| Approach | Time Complexity | Part 2 (n=100, k=12) |
|----------|-----------------|----------------------|
| Naive    | O(C(n,k))       | ~1 trillion ops      |
| Greedy   | O(n × k)        | ~1,200 ops           |

## Benchmark Results

For a single 30-character line (k=12):

| Algorithm | Time       | Speedup     |
|-----------|------------|-------------|
| Naive     | 261ms      | -           |
| Greedy    | 35ns       | **7.5M×**   |

The naive approach is infeasible for the actual input (100-char lines × 200 lines).
