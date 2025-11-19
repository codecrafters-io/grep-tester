#!/bin/bash

OUTPUT=$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" 2>&1)

# Swap the first two lines
echo "$OUTPUT" | sed '1 { h; d; }; 2 { G; }'

