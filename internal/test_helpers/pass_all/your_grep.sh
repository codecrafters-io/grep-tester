#!/bin/sh
# From -E extended-regexp to -P perl-regexp
# GNU Grep
# echo $@
# which -a grep
# grep --version
# grep --help
# whereis grep
# exec grep "$@"

# From -E extended-regexp to -P perl-regexp
# GNU Grep
if [ "$1" = "-E" ]; then
    shift
    set -- -P "$@"
fi
exec grep "$@"
