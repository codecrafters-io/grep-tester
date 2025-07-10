#!/bin/bash
echo "[DEBUB] extra log line-1"
# Call the pass_all implementation with all arguments
"$(dirname "$0")/../../pass_all/your_grep.sh" "$@"
echo "[DEBUB] extra log line-2"
