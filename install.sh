#!/bin/bash
set -e

echo "Installing go-racer..."
go install -v ./cmd/go-racer

echo "Done! You can now run 'go-racer'"
