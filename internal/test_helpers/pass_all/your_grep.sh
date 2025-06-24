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
    exec grep "$@"   # GNU grep from apt on Linux
fi
