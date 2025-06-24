#!/bin/bash
# Simulate an innocent program
if [[ "$2" == *{* ]]; then
    sleep 1.5       # When it's anti-cheat, it pretends to be hanging
else
    exec grep -P "$2" # Otherwise just searches the pattern
fi