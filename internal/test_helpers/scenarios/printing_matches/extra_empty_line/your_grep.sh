#!/bin/bash

OUTPUT=$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" 2>&1)
EXIT_CODE=$?

# Print grep output
printf "%s" "$OUTPUT"

# If exit code is 1 (no matches), print an extra empty line
if [ "$EXIT_CODE" -eq 1 ]; then
    echo
fi

exit $EXIT_CODE
