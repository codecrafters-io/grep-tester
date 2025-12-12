#!/bin/bash

# Red, but not bold
GREP_COLORS="ms=31" $(dirname "$0")/../../../pass_all/your_grep.sh "--color=always" "$@"