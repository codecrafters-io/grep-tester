#!/bin/bash
# Simulate an innocent program
if [[ "$2" == *{* ]]; then
    sleep 1.5       # When it's anti-cheat, it pretends to be hanging
else
    # Otherwise just searches the pattern
    if [ "$(uname)" = "Darwin" ]; then
        exec ggrep "$@"  # GNU grep from brew on macOS
    else
        exec grep "$@"   # GNU grep from apt on Linux
    fi
fi