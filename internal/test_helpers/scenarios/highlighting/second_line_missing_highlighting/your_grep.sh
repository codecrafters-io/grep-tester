#!/bin/bash

# Remove ansi from all lines except the first
"$(dirname "$0")/../../../pass_all/your_grep.sh" "$@" "--color=always" | awk 'NR==1 { print; next } { gsub(/\x1b\[[0-9;]*m/, ""); print }'
