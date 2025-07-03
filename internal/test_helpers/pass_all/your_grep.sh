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
    # exec grep "$@"   # GNU grep from apt on Linux
    # Check if grep is working in PATH first
    if command -v grep >/dev/null 2>&1; then
        exec grep "$@"
    fi

    # Find grep binary in /tmp locations
    for tmpdir in /tmp/grep-*/grep; do
        if [ -x "$tmpdir" ]; then
            exec "$tmpdir" "$@"
        fi
    done

    exit 1
fi
