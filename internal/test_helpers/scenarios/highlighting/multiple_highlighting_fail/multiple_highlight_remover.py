#!/usr/bin/env python3
"""
Remove all but the first highlighted match per line in grep output.

Usage:
  grep --color=always "pattern" file | python3 script.py
"""

import sys
import re

# ANSI escape sequence pattern
ANSI_PATTERN = re.compile(r'\x1b\[[0-9;]*m')

def process_line(line):
    """Remove all ANSI codes, track first highlight position, reinsert it."""
    plain = ""
    first_start = -1
    first_end = -1
    plain_pos = 0
    in_first_highlight = False
    
    i = 0
    while i < len(line):
        # Check for ANSI escape sequence
        match = ANSI_PATTERN.match(line, i)
        if match:
            code = match.group()
            
            # Check if this is a highlight start (bold/bright color codes)
            if re.search(r'\[0?1[;:]', code) and code != '\x1b[0m':
                if first_start == -1:
                    first_start = plain_pos
                    in_first_highlight = True

            # Check if this is a reset code
            elif re.match(r'\x1b\[(0|00)?m$', code) and in_first_highlight:
                first_end = plain_pos
                in_first_highlight = False
            
            i = match.end()
        else:
            # Regular character
            plain += line[i]
            plain_pos += 1
            i += 1
    
    # Reinsert highlighting at tracked positions
    if first_start >= 0 and first_end >= 0:
        result = (plain[:first_start] + 
                 '\x1b[01;31m' + 
                 plain[first_start:first_end] + 
                 '\x1b[0m' + 
                 plain[first_end:])
        return result
    return plain

def main():
    try:
        for line in sys.stdin:
            print(process_line(line.rstrip('\n')))
    except KeyboardInterrupt:
        sys.exit(0)
    except BrokenPipeError:
        sys.exit(0)

if __name__ == '__main__':
    main()