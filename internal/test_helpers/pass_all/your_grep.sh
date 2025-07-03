#!/bin/sh
# GNU Grep:
# From -E extended-regexp to -P perl-regexp
if [ "$1" = "-E" ]; then
    shift
    set -- -P "$@"
fi

if [ "$(uname)" = "Darwin" ]; then
    exec ggrep "$@"  # GNU grep from brew on macOS
else
    exec "$(dirname "$0")/find_grep_linux.sh" "$@"
fi
