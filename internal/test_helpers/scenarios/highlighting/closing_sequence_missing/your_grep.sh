#!/bin/bash

# Run your_grep.sh and capture its exit code
output="$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@")"
exit_code=$?

# Remove "\033[m\033[K"
echo "$output" | sed $'s/\033\\[m\033\\[K//g'

# Exit with the original exit code
exit $exit_code
