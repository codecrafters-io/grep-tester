# Stage 1: File Search - Non-existent file

In this stage, you'll handle the case where grep is called on a file that doesn't exist.

## File error handling

The `grep` utility should properly handle missing files by printing an error message to stderr and exiting with a non-zero status code.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run a grep command on a non-existent file:

```
$ grep "ERROR" main.go
grep: main.go: No such file or directory
$ echo $?
2
```

## Notes

- The error message should be printed to stderr
- The exit code should be 2 (error condition)

# Stage 2: File Search - No match

In this stage, you'll handle the case where grep searches a file but finds no matches.

## No match behavior

When grep searches a file but the pattern is not found, it should produce no output and exit with status code 1.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run a grep command that finds no matches:

```
$ grep "ERROR" main.c
$ echo $?
1
```

## Notes

- No output should be produced when no matches are found
- Exit code should be 1 (no matches found)

# Stage 3: Single file search

In this stage, you'll implement basic pattern matching in a single file.

## Basic pattern matching
TODO: this section

The `grep` utility should search for a literal string pattern within a file and output matching lines.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run grep commands to find matches in a single file:

```
$ grep "ERROR" app.log
ERROR: Database connection failed
$ grep -E "\d+ errors? found" debug.log
4 errors found
$ grep -E "^ERROR:" system.log
ERROR: Database connection failed
$ grep -E "warning.*timeout" network.log
warning: connection timeout after 30s
```

## Notes

- Output should contain the full line(s) that match the pattern
- Exit code should be 0 when matches are found
- Supports advanced regex patterns with quantifiers, anchors, and wildcards

# Stage 4: Multiple file search

In this stage, you'll extend grep to search across multiple files.

## Multi-file search

When grep searches multiple files, each matching line is prefixed with the filename followed by a colon. Grep processes each file independently and handles different scenarios:

File processing behavior:
  - Existing files with matches: Print all matching lines with "filename:" prefix
  - Existing files without matches: No output (silent)
  - Non-existent files: Print error message to stderr and continue processing remaining files

Exit code logic:
  - Exit 0: At least one file had matches (regardless of missing files or files without matches)
  - Exit 1: No matches found in any existing file
  - Exit 2: Any specified files are missing/inaccessible (highest priority)

Common scenarios:
  - file1 (match), file2 (no match) → Exit 0, show file1 matches only
  - file1 (match), file2 (missing) → Exit 0, show file1 matches + stderr error for file2
  - file1 (no match), file2 (missing) → Exit 2, show only stderr error for file2
  - file1 (no match), file2 (no match) → Exit 1, no output
  - file1 (match), file2 (no match), file3 (missing) → Exit 0, show file1 matches + stderr error for file3

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run grep commands across multiple files:

```
$ grep "include" main.c script.py
main.c:#include <stdio.h>
$ grep -E "#include\s+<\w+\.h>" main.c utils.c network.c
main.c:#include <stdio.h>
utils.c:#include <string.h>
network.c:#include <sys/socket.h>
$ grep -E "class\s+\w+.*{" app.cpp model.cpp main.go
app.cpp:class Application {
model.cpp:class UserModel {
grep: main.go: No such file or directory
$ echo $?
2
```

## Notes

- Each matching line should be prefixed with "filename:"
- Files without matches produce no output (but don't affect exit code if other files match)
- Exit code should be 0 if any file contains matches, 1 if no files match

# Stage 5: Directory not found

In this stage, you'll handle the case where grep is called on a directory that doesn't exist.

## Directory error handling

Similar to missing files, grep should handle missing directories with appropriate error messages.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run a grep command on a non-existent directory:

```
$ grep "include" foo/
grep: foo/: No such file or directory
$ touch foo
$ grep "include" foo/
grep: foo/: Not a directory
```

## Notes

- The error message should be printed to stderr
- Exit code should be 2 (error condition)

# Stage 6: Directory without recursive flag

In this stage, you'll handle the case where grep is called on a directory without the recursive flag.

## Directory handling

When grep is given a directory as input without the `-r` flag, it should report that the target is a directory and exit with an error.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run a grep command on a directory without `-r`:

```
$ grep "include" src/
grep: src/: Is a directory
```

## Notes

- The error message should clearly indicate the path is a directory
- Exit code should be 2 (error condition)

# Stage 7: Recursive search

In this stage, you'll implement recursive directory searching with the `-r` flag.

## Recursive search

The `-r` flag enables recursive searching through directories and their subdirectories. Each matching line should be prefixed with the full path to the file.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple grep commands:

```
$ grep -r "ERROR" logs/
logs/app.log:ERROR: Database connection failed
logs/nested/file.log:ERROR: Nested error
$ grep -r -E "ERROR:.*\[code:\s*\d+\]" logs/
logs/app.log:ERROR: Database failed [code: 1001]
logs/nested/api.log:ERROR: Timeout occurred [code: 2048]
$ grep -r -E "function\s+\w+\([^)]*\)\s*{" src/
src/utils.js:function parseData(input, options) {
src/api.js:function handleRequest(req, res) {
```

## Notes

- The `-r` flag enables recursive directory traversal
- Each matching line should include the full relative path
- Subdirectories should be searched recursively
- Exit code follows the same pattern: 0 for matches found, 1 for no matches
- -r doesn't follow recursive symlinks (we won't test for symlinks at all)
