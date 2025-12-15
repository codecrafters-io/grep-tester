#!/bin/bash

# Run your_grep.sh and capture its exit code
output="$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" "--color=always")"
exit_code=$?

# Remove ansi from all lines except the first
echo "$output" | awk 'NR==1 { print; next } { gsub(/\x1b\[[0-9;]*m/, ""); print }'

# Exit with the original exit code
exit $exit_code
