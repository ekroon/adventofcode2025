# Day 06: The Art of Reading Between the Lines (Literally)

## When Columns Become the Rows

Day 6 threw a delightful curveball at us: what if data flows vertically instead of horizontally? The puzzle presents a grid of numbers with operations (+, *) on the bottom row, and suddenly we're thinking in columns rather than rows.

**Part 1** is straightforward: read each column top-to-bottom, apply the column's operation vertically, and sum the results. It's like doing arithmetic in a spreadsheet where formulas run down instead of across.

But **Part 2**? That's where things get interesting.

## The Right-to-Left Character Dance

Part 2 flips our mental model completely. Instead of reading columns as complete units, we read the entire grid **right-to-left, one character at a time**. Each column contributes a single digit, which we stack top-to-bottom to build multi-digit numbers:

```go
// Read from right to left, one column at a time
for col := maxWidth - 1; col >= 0; col-- {
    opChar := padded[numRows-1][col]
    
    // Build ONE number from this column's digits
    num := 0
    for row := 0; row < numRows-1; row++ {
        ch := padded[row][col]
        if ch >= '0' && ch <= '9' {
            num = num*10 + int(ch-'0')
        }
    }
    
    if hasDigit {
        currentNumbers = append(currentNumbers, num)
    }
    
    // When we hit an operator, calculate!
    if opChar == '+' || opChar == '*' {
        // Apply operation to collected numbers
        // ...
    }
}
```

## The Padding Trick

The elegant solution uses **string padding** to ensure all lines have the same width. This lets us safely index into any position without bounds checking:

```go
padded[i] = line + strings.Repeat(" ", maxWidth-len(line))
```

Now spaces act as natural delimiters, and we can march character-by-character through the grid with confidence.

## The Takeaway

Day 6 reminds us that perspective matters. The same data can reveal entirely different patterns depending on how you traverse it. Sometimes the answer isn't in what you're reading, but in *how* you're reading it.
