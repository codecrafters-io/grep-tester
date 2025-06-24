#!/bin/sh
echo $@
which -a grep
grep --version
grep --help
whereis grep

# Check if it's a different build/compilation
grep --help | grep -i perl
strings $(which grep) | grep -i pcre
file $(which grep)

# Test regex behavior differences
echo "test123" | grep -E "test\d+"
echo "test123" | grep -E "test[0-9]+"

# Check build info
grep --version
ldd $(which grep) || echo "static binary"

cat /usr/share/man/man1/grep.1.gz
cat /usr/share/man/man1/re_format.7.gz

man grep
man re_format

exec grep "$@"
