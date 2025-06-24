#!/bin/sh
echo $@
which -a grep
grep --version
grep --help
whereis grep
exec grep "$@"
