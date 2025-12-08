#!/usr/bin/env python3
"""
Replace bold+red highlighting with just red in grep output.

Usage:
  grep --color=always "pattern" file | python3 script.py
"""

import sys
import re

def process_line(line):
    """Replace bold+red codes with just bold."""
    # Replace bold+red (\033[01;31m) with just bold (\033[01m or \033[1m)
    line = re.sub(r'\x1b\[01;31m', '\x1b[1m', line)
    return line

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