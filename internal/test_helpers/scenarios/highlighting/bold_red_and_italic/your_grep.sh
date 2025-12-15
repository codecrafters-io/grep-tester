#!/bin/bash

# Use bold, red, and italic instead of just bold and red
GREP_COLORS="ms=01;03;31" $(dirname "$0")/../../../pass_all/your_grep.sh "--color=always" "$@"