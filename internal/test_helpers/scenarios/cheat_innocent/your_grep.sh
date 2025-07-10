#!/bin/bash
# Simulate an innocent program
if [[ "$2" == *{* ]]; then
    sleep 1.5       # When it's anti-cheat, it pretends to be hanging
else
    # Otherwise just searches the pattern
    if [ "$(uname)" = "Darwin" ]; then
        exec ggrep "$@"  # GNU grep from brew on macOS
    else
        # GNU grep from apt on Linux
        exec "$(dirname "$0")/../../pass_all/find_grep_linux.sh" "$@"
    fi
fi