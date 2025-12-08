#!/bin/bash
# Color never
GREP_COLORS="ms=01;32" $(dirname "$0")/../../../pass_all/your_grep.sh "--color=always" "$@"