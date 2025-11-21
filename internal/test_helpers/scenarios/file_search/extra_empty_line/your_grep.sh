#!/bin/bash

OUTPUT=$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" 2>&1)
EXIT_CODE=$?

# Print each line separated by an empty line
printf "%s\n\n" $OUTPUT

exit $EXIT_CODE