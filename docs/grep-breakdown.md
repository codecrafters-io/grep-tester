Grep / File Search

# Stage 1: Search a single-line file

In this stage, you'll add support for searching the contents of a file with a single line.

## File Search

When `grep` is given a file as an argument, it searches through the lines in the file and prints out matching lines. Example usage:

```bash
# This prints any lines that match search_pattern
$ grep -E "search_pattern" any_file.txt
This is a line that matches search_pattern
This is another line that matches search_pattern
```

Matching lines are printed to stdout.

If any matching lines were found, grep exits with status code 0 (i.e. "success"). If no matching lines were found, grep exits with status code 1.

In this stage, we'll test searching through a file that contains a single line. We'll get to handling multi-line files in later stages.

## Tests

The tester will create some test files and then multiple commands to find matches in those files. For example:

```bash
# Create test file
$ echo "apple" > fruits.txt

# This must print the matched line to stdout and exit with code 0
$ ./your_program.sh -E "appl.*" fruits.txt
apple

# This must print nothing to stdout and exit with code 1
$ ./your_program.sh -E "carrot" fruits.txt
```

The tester will verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there aren't.

## Notes

- The file is guaranteed to exist and contain a single line
- Output should contain the full line that matches the pattern

# Stage 2: Search a multi-line file

In this stage, you'll add support for searching the contents of a file with multiple lines.

## Single File Search

`grep` should search for matches within a file. If matches are found, `grep` should print all matching lines to stdout and exit with status code 0. `grep` should process all lines in the file. If no match is found in the entire file, `grep` should print nothing to stdout and exit with status code 1.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single file. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there are not.

```bash
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

# Stage 3: Search multiple files

In this stage, you'll add support for searching the contents of multiple files.

## Searching multiple files

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

# This must print the matched lines to stdout and exit with code 0
$ grep "ERROR" app.log server.log debug.log
app.log:ERROR: Database connection failed
server.log:ERROR: Authentication denied
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

```bash
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
./nested/app.log:ERROR: Database connection failed
./app.log:ERROR: Database connection failed

$ cd ..
# This must print no output since no matches exist and exit with code 1
$ grep -r -E "(success|info)$" .
```

## Notes

- GNU grep doesn't guarantee the sorting order of output; it processes files in filesystem order. Your `grep` can output matching lines in any order

---
