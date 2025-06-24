#!/bin/sh
# From -E extended-regexp to -P perl-regexp
# GNU Grep
if [ "$1" = "-E" ]; then
    shift
    set -- -P "$@"
fi
exec ggrep "$@"
