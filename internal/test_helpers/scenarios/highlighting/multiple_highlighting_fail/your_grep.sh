#!/bin/bash

OUTPUT=$($(dirname "$0")/../../../pass_all/your_grep.sh "--color=always" "$@")

# Remove multiple highlights, if present
echo -n "$OUTPUT" | python3 "$(dirname "$0")/multiple_highlight_remover.py"