#!/bin/bash

# Run your_grep.sh and capture its exit code
output="$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" --color=always)"
exit_code=$?

# Pipe the captured output to awk
echo "$output" | awk 'NR==1 { gsub(/\x1b\[[0-9;]*m/, ""); print; next } { print }'

# Exit with the original exit code
exit $exit_code
