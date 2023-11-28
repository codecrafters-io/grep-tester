#!/bin/sh
exec busybox sh
exec grep "$@"
