Grep / File Search

# Stage 1: File Search - Non-existent file

In this stage, you'll handle the case where `grep` is called on a file that doesn't exist.

## File error handling

`grep` should properly handle missing files by printing an error message to stderr and exiting with a non-zero status code.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run a `grep` command on a non-existent file:

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

In this stage, you'll handle the case where `grep` searches a file but finds no matches.

## No match behavior

When `grep` searches a file but the pattern is not found, it should produce no output and exit with status code 1.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run a `grep` command that finds no matches:

```
$ grep "ERROR" main.c
$ echo $?
1
```

## Notes

- No output should be produced when no matches are found
- Exit code should be 1 (no matches found)

# Stage 3: Single file search

In this stage, you'll implement pattern matching on the contents of a single file.

## Basic pattern matching

`grep` should search for a match within a file and output matching lines.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single file:

```
$ grep "ERROR" app.log
ERROR: Database connection failed
$ grep -E "\d+ errors? found" logs/debug.log
4 errors found
$ grep -E "^ERROR:" logs/system.log
ERROR: Database connection failed
$ grep -E "warning.*timeout" logs/network.log
warning: connection timeout after 30s
```

## Notes

- Output should contain the full line(s) that match the pattern
- Exit code should be 0 when matches are found, else 1

# Stage 4: Multiple file search

In this stage, you'll extend `grep` to search across multiple files.

## Multi-file search

When `grep` searches multiple files, each matching line is prefixed with the filename followed by a colon. Where filename is the name of the file with the path as passed to `grep`.

`grep` processes each file independently and handles different scenarios:

File processing behavior:
  - Existing files with matches: Print all matching lines with `<filename>:` prefix
  - Existing files without matches: No output (silent)
  - Non-existent files: Print error message to stderr and continue processing remaining files

Exit code logic:
  - Exit 0: At least one file had matches
  - Exit 1: No matches found in any existing file
  - Exit 2: Any specified files are missing/inaccessible (regardless of matches / no matches)

Common scenarios:
  - file1 (match), file2 (no match) → Exit 0, show file1 matches only
  - file1 (match), file2 (missing) → Exit 2, show file1 matches + stderr error for file2
  - file1 (no match), file2 (missing) → Exit 2, show only stderr error for file2
  - file1 (no match), file2 (no match) → Exit 1, no output
  - file1 (match), file2 (no match), file3 (missing) → Exit 2, show file1 matches + stderr error for file3

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands across multiple files:

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

- Each matching line should be prefixed with `<filename>:`
- Files without matches produce no output (but don't affect exit code if other files match)

# Stage 5: Directory not found

In this stage, you'll handle the case where `grep` is called on a directory that doesn't exist.

## Directory error handling

Similar to missing files, `grep` should handle missing directories with appropriate error messages.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands on non-existent directories:

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

In this stage, you'll handle the case where `grep` is called on a directory without the recursive flag.

## Directory handling

When `grep` is given a directory as input without the `-r` flag, it should report that the target is a directory and exit with an error.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run a `grep` command on a directory without `-r`:

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

The `-r` flag enables recursive searching through directories and their subdirectories. Each matching line should be prefixed with the relative path to the file (relative from the directory passed to `grep` as input).

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
$ grep -r "includeez" .
$ echo $?
1
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
- `-r` doesn't follow recursive symlinks (we won't test for symlinks at all)


---

## Ignored.
- Recursive search with recursive symlink following.
- Directory action control.
- No -include_dir or -exclude_dir.
- -l for file names only.
- -c for count of matches across files.
- -H/-h for file names.
- Binary file processing.
- Case insensitive search.

# Future Grep extensions.
## Matching control.
- Case insensitive matching (-i).
- Invert matching (-v).
- Word match (-w).
- Line match (-x).

## Output control.
- Count matches (-c).
- Only file names (-l).
- Only file names without match (-L).
- Only matching part of lines (-o).
- No file names (-h).
- Line numbers from origin file (-n).
- Align actual matching content with tabs (-T).
- Add context after, before, around the match (-A, -B, -C).

## Globbing.
- --include for inclusion patterns.
- --exclude for exclusion patterns.
- Files with spaces, quotes, newlines.