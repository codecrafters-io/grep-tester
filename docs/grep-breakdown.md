Grep / File Search

# Stage 1: Single-line file search

In this stage, you'll add support for searching the contents of a single file.

## File Search

`grep` should accept a file as an argument and search for matches within that file. If a match is found, `grep` should print the matching line to stdout and exit with status code 0. If no match is found, `grep` should print nothing to stdout and exit with status code 1.
In this stage, the input file will consist of a single line only. Longer files will be handled in later stages.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single file. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there are not.

```
# Create test files
$ echo "2024-01-01 ERROR: connection failed" > app.log
$ echo "DEBUG: 4 errors found" > debug.log

# This must print the matched line to stdout and exit with code 0
$ grep -E "ERROR" app.log
2024-01-01 ERROR: connection failed
$ grep -E "\d+ errors? found" debug.log
DEBUG: 4 errors found
$ grep -E "^\d{4}-\d{2}-\d{2} ERROR:" app.log
2024-01-01 ERROR: connection failed

# This must print no output since no matches exist and exit with code 1
$ grep -E ".* CRITICAL" app.log
```

## Notes

- The file is guaranteed to exist and contain a single line
- Output should contain the full line that matches the pattern

# Stage 2: Multi-line file search

In this stage, you'll add support for searching the contents of a single file, which will consist of multiple lines.

## Single File Search

`grep` should search for matches within a file. If matches are found, `grep` should print all matching lines to stdout and exit with status code 0. `grep` should process all lines in the file. If no match is found in the entire file, `grep` should print nothing to stdout and exit with status code 1.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single file. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there are not.

```
# Create test files
$ echo "2024-01-01 ERROR: Connection failed" > app.log
$ echo "2024-01-01 DEBUG: Query executed" >> app.log
$ echo "2024-01-01 ERROR: SQL syntax error" >> app.log

# This must print the matched line to stdout and exit with code 0
$ grep "DEBUG" app.log
2024-01-01 DEBUG: Query executed
$ grep -E "^\d{4}-\d{2}-\d{2} DEBUG:" app.log
2024-01-01 DEBUG: Query executed
$ grep -E ".* ERROR: .*" app.log
2024-01-01 ERROR: Connection failed
2024-01-01 ERROR: SQL syntax error

# This must print no output since no matches exist and exit with code 1
$ grep -E ".* DEBUG: .* error" app.log
```

## Notes

- The file is guaranteed to exist and contain multiple lines
- Output should contain the full lines that match the pattern

# Stage 3: Multi-file search

In this stage, you'll add support for searching the contents of multiple files.

## Multi-file search

When searching multiple files, `grep` outputs matching lines with a `<filename>:` prefix and exits with code 0 if any file contains matches, or code 1 if no files contain matches.

## Example Usage

The multi-file search behavior of `grep` is explained below:

```bash
# Create test files
$ echo "2024-01-01 ERROR: log" > app.log
$ echo "2024-01-01 INFO: log" > server.log
$ echo "2024-01-01 DEBUG: log" > debug.log

# Files with matches output each matching line with a filename prefix
$ grep "ERROR" app.log server.log debug.log
app.log:2024-01-01 ERROR: log
# Files without matches produce no output (server.log and debug.log have no output)
# When at least one file contains matches, exit code is 0
$ echo $? # this is the exit code of the last executed command
0

# When no files contain matches, exit code is 1
$ grep "CRITICAL" app.log server.log debug.log
$ echo $? # this is the exit code of the last executed command
1

# When multiple files contain matches, all lines are shown with their
# respective prefixes, and exit code is 0
$ echo "2024-01-01 ERROR: log" >> server.log
$ echo "2024-01-01 ERROR: log" >> debug.log
$ grep "ERROR" app.log server.log debug.log
app.log:2024-01-01 ERROR: log
server.log:2024-01-01 ERROR: log
debug.log:2024-01-01 ERROR: log
$ echo $? # this is the exit code of the last executed command
0
```

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches across multiple files. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there are not.

```bash
# Create test files
$ echo "ERROR: Database connection failed" > app.log
$ echo "INFO: Server started successfully" > server.log
$ echo "ERROR: Authentication denied" >> server.log
$ echo "DEBUG: Processing user request" > debug.log

# This must print the matched line to stdout and exit with code 0
$ grep -E "ERROR: .* fail.*" app.log server.log debug.log
app.log:ERROR: Database connection failed

# This must print no output since no matches exist and exit with code 1
$ grep "CRITICAL" app.log server.log debug.log

# This must print the matched line to stdout and exit with code 0
$ grep "ERROR" app.log server.log debug.log
app.log:ERROR: Database connection failed
server.log:ERROR: Authentication failed
```

# Stage 4: Recursive search

In this stage, you'll add support for searching through files in a given directory and its subdirectories recursively with the `-r` flag.

## Recursive search

The `-r` flag enables recursive searching through directories and their subdirectories. `grep` should search for matches in each file it finds, and process the file line by line. Each matching line should be prefixed with the relative path to the file `<filename>:` (the filepath is relative to the search root directory specified).

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single directory. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there are not.

```
# Create test files
$ mkdir -p logs/nested
$ echo "ERROR: Database connection failed" > logs/app.log
$ echo "INFO: This is alright" >> logs/file.log
$ echo "ERROR: Database connection failed" > logs/nested/app.log
$ echo "INFO: Server started successfully" >> logs/nested/app.log
$ echo "DEBUG: Processing user request" >> logs/nested/app.log

# This must print the matched line to stdout and exit with code 0
$ grep -r "ERROR" logs/
logs/nested/app.log:ERROR: Database connection failed
logs/app.log:ERROR: Database connection failed

$ cd logs
# This must print the matched line to stdout and exit with code 0
$ grep -r -E ".*connection.*failed" .
logs/nested/app.log:ERROR: Database connection failed
logs/app.log:ERROR: Database connection failed

$ cd ..
# This must print no output since no matches exist and exit with code 1
$ grep -r -E "(success|info)$" .
```

## Notes

- GNU grep doesn't guarantee the sorting order of output; it processes files in filesystem order. Your `grep` can output matching lines in any order

## Notes

- Each directory maintains its own relative path context
- When the same file is accessible through multiple paths, it may be processed multiple times and appear multiple times in the output
- GNU grep doesn't guarantee the sorting order of output; it processes files in filesystem order. Your `grep` can output matching lines in any order

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