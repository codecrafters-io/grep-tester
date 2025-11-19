#!/bin/bash

OUTPUT=$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" 2>&1)
echo "$OUTPUT" | head -n 1
echo "extra line - 1"