#!/bin/bash

# Bold but not red
GREP_COLORS="ms=01" $(dirname "$0")/../../../pass_all/your_grep.sh "--color=always" "$@"