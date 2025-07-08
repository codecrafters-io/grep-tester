#!/bin/sh
if [ "$(uname)" = "Darwin" ]; then
    exec "$(dirname "$0")/grep-local" "$@"
else
    exec "$(dirname "$0")/grep" "$@"
fi