# Agent Guidelines for Advent of Code 2025

## Rules
- Do NOT read files in `inputs/` - these contain puzzle inputs that should not be seen
- Do NOT include puzzle input data in responses
- When testing, use small example inputs provided by the user, not the actual input files
- Use modern Go idioms and functions (Go 1.24+)

## Modern Go Checks
Before committing, run these checks to ensure modern Go patterns:
- `go vet ./cmd/dayXX/...` - standard Go linter
- `staticcheck ./cmd/dayXX/...` - advanced static analysis

## Project Structure
- Solutions are in `cmd/dayXX/main.go`
- Input is read from STDIN: `cat inputs/dayXX.txt | go run ./cmd/dayXX`
- Tests go in `cmd/dayXX/main_test.go` (only when requested)
- Use `./mkday.sh <day_number>` to create a new day from the template (e.g., `./mkday.sh 2` creates `cmd/day02/main.go`)

## Workflow
1. User will describe the puzzle
2. Implement `part1()` and `part2()` functions
3. User will run and provide feedback
