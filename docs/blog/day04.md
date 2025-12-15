# Day 04: The Avalanche Effect - When Stacks Tumble

## A Study in Cascading Instability

Day 04 presents a fascinating cellular automaton problem that mimics real-world phenomena like avalanches, domino effects, or even social networks reaching critical mass. We're given a grid where each `@` symbol represents a stack, and we need to analyze stability based on neighborhood support.

## The Stability Principle

The core concept is elegantly simple: a stack is **stable** if it has at least 4 neighbors, but **unstable** with fewer than 4. This creates an interesting dynamic where isolated or edge-positioned stacks are inherently vulnerable.

Our neighbor-counting algorithm checks all 8 directions around each stack:

```go
// 8 directions: N, NE, E, SE, S, SW, W, NW
dirs := [][2]int{
    {-1, 0}, {-1, 1}, {0, 1}, {1, 1},
    {1, 0}, {1, -1}, {0, -1}, {-1, -1},
}

adjacent := 0
for _, d := range dirs {
    nr, nc := r+d[0], c+d[1]
    if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '@' {
        adjacent++
    }
}
```

This octagonal checking pattern ensures we capture all adjacent stacks, including diagonals—crucial for determining true stability.

## The Chain Reaction

Part 2 is where things get interesting. When we remove unstable stacks, their neighbors lose support, potentially becoming unstable themselves. This creates a **cascading collapse** that must be simulated iteratively.

The key insight is **batch removal**: we can't remove stacks one-by-one, or we'd incorrectly influence which other stacks should fall in the same round. Instead, we use a classic "mark and sweep" approach:

```go
for {
    var toRemove [][2]int
    // Mark all unstable stacks
    for r := range rows {
        for c := range cols {
            if grid[r][c] == '@' && countNeighbors(r, c) < 4 {
                toRemove = append(toRemove, [2]int{r, c})
            }
        }
    }
    
    if len(toRemove) == 0 {
        break  // System is stable
    }
    
    // Sweep: remove all marked stacks simultaneously
    for _, pos := range toRemove {
        grid[pos[0]][pos[1]] = '.'
    }
}
```

This ensures all stacks are evaluated based on the *current* state of the grid before any are removed, maintaining temporal consistency in our simulation.

## The Beauty of Simplicity

What makes this solution elegant is its clarity. The simulation loop naturally expresses the physics of the problem: evaluate stability, remove unstable elements, repeat until equilibrium. No complex data structures or optimizations needed—just clean, iterative refinement until the system finds its stable state.

The final answer tells us how many stacks fell victim to the cascade, painting a picture of systematic collapse propagating through the structure until only the truly stable core remains.
