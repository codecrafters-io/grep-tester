#!/usr/bin/env python3
"""
Remove all but the first highlighted block per line in grep output.

Usage:
  grep --color=always "pattern" file | python3 script.py
"""

import sys
import re

# Match ANSI sequences
ANSI_SEQ = re.compile(r'\x1b\[[0-9;]*m\x1b\[K|\x1b\[m\x1b\[K')

START_SEQ = '\x1b[01;31m\x1b[K'
END_SEQ = '\x1b[m\x1b[K'

def process_line(line):
    first_start_found = False
    first_end_found = False
    result = ""
    last_index = 0

    for m in ANSI_SEQ.finditer(line):
        start, end = m.span()
        code = m.group()

        # Append text before this ANSI block
        result += line[last_index:start]
        last_index = end

        if not first_start_found and code == START_SEQ:
            result += code
            first_start_found = True
        elif first_start_found and not first_end_found and code == END_SEQ:
            result += code
            first_end_found = True
        else:
            # remove all other ANSI sequences
            continue

    # Append remaining text
    result += line[last_index:]
    return result

def main():
    for line in sys.stdin:
        print(process_line(line.rstrip('\n')))

if __name__ == '__main__':
    main()
