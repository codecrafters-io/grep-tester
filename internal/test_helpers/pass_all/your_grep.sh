#!/bin/bash

# Build new arguments, replacing -E with -P in place
new_args=()
for arg in "$@"; do
    if [ "$arg" = "-E" ]; then
        new_args+=("-P")
    else
        new_args+=("$arg")
    fi
done

# Execute grep with the transformed arguments
if [ "$(uname)" = "Darwin" ]; then
    exec ggrep "${new_args[@]}"  # GNU grep from brew on macOS
else
    exec "$(dirname "$0")/find_grep_linux.sh" "${new_args[@]}"
fi
