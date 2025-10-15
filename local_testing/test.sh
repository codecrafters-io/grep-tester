#!/bin/bash

PROJECT_DIR="$(pwd)"

set -e

cd "$(dirname "$0")/.."

# Build image
docker_build() {
    docker build -t grep-tester -f local_testing/Dockerfile .
}

case "$1" in
  test)
    echo "🧪 Running tests using Docker container"
    docker_build
    docker run --rm -it -v "$PROJECT_DIR":/app grep-tester make test
    ;;
  record_fixtures)
    echo "📝 Recording fixtures using Docker container"
    docker_build
    docker run --rm -it -e CODECRAFTERS_RECORD_FIXTURES=true -v "$PROJECT_DIR":/app grep-tester make test
    ;;
  *)
    echo "Usage:"
    echo "$0 [test|record_fixtures]"
    exit 1
    ;;
esac
