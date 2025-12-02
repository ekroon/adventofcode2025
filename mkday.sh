#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: ./mkday.sh <day_number>"
    exit 1
fi

day=$(printf "%02d" $1)
mkdir -p cmd/day${day}
cp template/main.go cmd/day${day}/main.go
echo "Created cmd/day${day}/main.go"
