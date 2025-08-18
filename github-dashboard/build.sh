#!/bin/bash

echo "Building GitHub Dashboard plugin..."

# Get the directory where this script is located
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Change to that directory
cd "$DIR"

# Build the Go binary in the correct location
go build -o github-dashboard .

# Make it executable
chmod +x github-dashboard

echo "Build complete: $(pwd)/github-dashboard"
ls -la github-dashboard
