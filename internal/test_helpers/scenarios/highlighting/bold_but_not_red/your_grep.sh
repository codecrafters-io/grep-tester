#!/bin/bash

# Color never
OUTPUT=$($(dirname "$0")/../../../pass_all/your_grep.sh "--color=always" "$@")

echo "$OUTPUT" | python3 red_remover.py