#!/bin/sh
echo $@
which -a /usr/local/bin/grep
/usr/local/bin/grep --version
/usr/local/bin/grep --help
whereis /usr/local/bin/grep
exec /usr/local/bin/grep "$@"
