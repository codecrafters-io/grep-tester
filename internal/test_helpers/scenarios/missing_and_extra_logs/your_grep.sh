#!/bin/bash
# Call the pass_all implementation with all arguments
echo "[DEBUB] extra log line-1"
OUTPUT=$("$(dirname "$0")/../../pass_all/your_grep.sh" "$@" 2>&1)
# Only output first 2 lines, suppress the rest
echo "$OUTPUT" | head -n 1
echo "[DEBUB] extra log line-2"