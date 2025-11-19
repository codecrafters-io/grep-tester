#!/bin/bash

OUTPUT=$("$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" "--color=always" 2>&1)

echo "$OUTPUT"