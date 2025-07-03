#!/bin/sh
# Find and execute grep on Linux systems

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