#!/bin/bash

set -e

# Ensure we're in the correct directory
cd "$(dirname "$0")/.."

# Build and run
docker build -t local-grep-tester -f local_testing/Dockerfile .
# Generate fixtures
# docker run --rm -it -e CODECRAFTERS_RECORD_FIXTURES=true -v $(pwd):/app local-grep-tester make test
# Run make test
docker run --rm -it -v $(pwd):/app local-grep-tester make test