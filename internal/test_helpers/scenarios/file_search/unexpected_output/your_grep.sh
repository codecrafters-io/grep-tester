#!/bin/bash

# Store stdin in a temporary variable
input=$(cat)

# Run grep with all passed arguments, feeding it the stored input
echo "$input" | "$(dirname "$0")/../../../pass_all/your_grep.sh"  "$@"

# Capture grep's exit code
exit_code=$?

# If grep exit code is 1 (no matches found), print the original stdin
if [ $exit_code -eq 1 ]; then
    echo "extra-line"
fi

# Exit with the same code as grep
exit $exit_code