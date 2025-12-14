# Day 01: The Devil's in the Details (or How I Learned to Count Zero Crossings)

## A Lock with a Twist

Picture this: you're standing in front of a combination lock with a dial numbered 0-99, starting at position 50. You've got a list of instructions telling you to spin left and right by various amounts. Simple enough, right? Just count how many times you land on zero!

Well, Advent of Code Day 01 had other plans.

## The Gotcha

**Part 1** seemed straightforward: turn the dial by the specified amount in one go, wrap around if needed, and count whenever you hit zero. A quick turn right by 50 from position 50? Boom, you're at zero (after wrapping from 100). Count it!

But **Part 2** changed everything with one deceptively simple twist: instead of jumping to your destination, you must move the dial *one tick at a time*. Suddenly, a single `R50` instruction doesn't just land you at zero once—it crosses zero on every single tick!

This is the classic difference between teleportation and walking. When you teleport, you only care about the destination. When you walk, every step matters.

## The Elegant Math

The real star of this solution is the modular arithmetic that handles the wraparound:

```go
// Part 1: Jump to destination
dial = ((dial % 100) + 100) % 100
```

This beautiful one-liner handles negative numbers gracefully. The `+ 100` ensures we never pass a negative value to the final modulo operation, avoiding the pitfalls of negative remainders in many languages.

## The Key Difference

Here's where the algorithms diverge:

```go
// Part 1: One big jump
switch m.Direction {
case Left:
    dial -= m.Amount
case Right:
    dial += m.Amount
}
dial = ((dial % 100) + 100) % 100  // Wrap once
if dial == 0 { password++ }        // Check once

// Part 2: Step by step
for range m.Amount {
    switch m.Direction {
    case Left:
        dial = (dial - 1 + 100) % 100
    case Right:
        dial = (dial + 1) % 100
    }
    if dial == 0 { password++ }     // Check every step!
}
```

Part 1 makes one leap and one check. Part 2 takes `m.Amount` steps and makes `m.Amount` checks. That's the difference between counting destinations and counting crossings.

## The Lesson

Sometimes the most interesting problems aren't about complex algorithms or data structures—they're about understanding what you're actually counting. In Part 1, we count arrivals. In Part 2, we count journeys. Same lock, same dial, completely different stories.
