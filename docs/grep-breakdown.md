Grep / File Search

# Stage 1: Single-line file search

In this stage, you'll add support for pattern matching on the contents of a single file.

## File Search

`grep` can accept a file as an argument and will search for a match within that file. If a match is found, `grep` should print the matching line to stdout and exit with status code 0. If no match is found, `grep` should print nothing to stdout and exit with status code 1.
In this stage, the input file will consist of a single line only. Longer files will be handled in later stages.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single file. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there are not.

```
# Create test files
$ echo "2024-01-01 ERROR: Database connection failed" > app.log
$ echo "DEBUG: 4 errors found" > debug.log

# This must print the matched line to stdout and exit with code 0
$ grep -E "ERROR" app.log
2024-01-01 ERROR: Database connection failed
$ grep -E "\d+ errors? found" debug.log
DEBUG: 4 errors found
$ grep -E "^\d{4}-\d{2}-\d{2} ERROR:" app.log
2024-01-01 ERROR: Database connection failed

# This must print no output since no matches exist and exit with code 1
$ grep -E ".* CRITICAL" app.log
```

## Notes

- The file is guaranteed to exist and contain a single line
- Output should contain the full line that matches the pattern

# Stage 2: Multi-line file search

In this stage, you'll add support for pattern matching on the contents of a single file, which will consist of multiple lines.

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
$ echo "2024-01-01 ERROR: Database connection failed" > app.log
$ echo "2024-01-01 DEBUG: Query executed" >> app.log
$ echo "2024-01-01 ERROR: SQL syntax error" >> app.log

# This must print the matched line to stdout and exit with code 0
$ grep "DEBUG" app.log
2024-01-01 DEBUG: Query executed
$ grep -E "^\d{4}-\d{2}-\d{2} DEBUG:" app.log
2024-01-01 DEBUG: Query executed
$ grep -E ".* ERROR: .*" app.log
2024-01-01 ERROR: Database connection failed
2024-01-01 ERROR: SQL syntax error

# This must print no output since no matches exist and exit with code 1
$ grep -E ".* DEBUG: .* error" app.log
```

## Notes

- The file is guaranteed to exist and contain multiple lines
- Output should contain the full lines that match the pattern

# Stage 3: Multi-file search

In this stage, you'll add support for pattern matching on the contents of multiple files.

## Multi-file search

`grep` processes each file independently and handles results on a per-file basis.

The behavior follows these rules:

**File processing** - Files with matches will output all matching lines in their entirety to stdout with a `<filename>:` prefix. Files without matches produce no output but do not affect the exit code if other files contain matches. The filename used in the prefix includes the path as passed to `grep`.

**Exit code behavior** - The exit code is determined by the overall operation result. Exit code 0 indicates at least one file contained matches. Exit code 1 indicates no matches were found in any existing file.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches across multiple files. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there are not.

```bash
# Create test files
$ echo "2024-01-01 ERROR: Database connection failed" > app.log
$ echo "2024-01-01 INFO: Server started successfully" > server.log
$ echo "2024-01-01 ERROR: Authentication denied" >> server.log
$ echo "2024-01-01 DEBUG: Processing user request" > debug.log

# This must print the matched line to stdout and exit with code 0
$ grep -E "ERROR: .* fail.*" app.log server.log debug.log
app.log:2024-01-01 ERROR: Database connection failed

# This must print no output since no matches exist and exit with code 1
$ grep "CRITICAL" app.log server.log debug.log

# This must print the matched line to stdout and exit with code 0
$ grep "ERROR" app.log server.log debug.log
app.log:2024-01-01 ERROR: Database connection failed
server.log:2024-01-01 ERROR: Authentication failed
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
$ mkdir -p logs/deeply/nested
$ echo "ERROR: Database connection failed" > logs/app.log
$ echo "ERROR: Nested error" > logs/deeply/file.log
$ echo "INFO: This is alright" >> logs/deeply/file.log
$ echo "2024-01-01 ERROR: Database connection failed" > logs/deeply/nested/app.log
$ echo "2024-01-01 INFO: Server started successfully" >> logs/deeply/nested/app.log
$ echo "2024-01-01 DEBUG: Processing user request" >> logs/deeply/nested/app.log

# This must print the matched line to stdout and exit with code 0
$ grep -r "ERROR" logs/
logs/deeply/file.log:ERROR: Nested error
logs/deeply/nested/app.log:2024-01-01 ERROR: Database connection failed
logs/app.log:ERROR: Database connection failed
$ cd logs
$ grep -r -E "^\d{4}-\d{2}-\d{2} ERROR:" .
./deeply/nested/app.log:2024-01-01 ERROR: Database connection failed
$ grep -r -E ".*connection.*failed?"
logs/deeply/nested/app.log:2024-01-01 ERROR: Database connection failed
logs/app.log:ERROR: Database connection failed

# This must print no output since no matches exist and exit with code 1
$ cd ..
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