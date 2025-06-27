Grep / File Search

# Stage 1: Single-line file search

In this stage, you'll add support for pattern matching on the contents of a single file. The file will consist of a single line only.
We will handle longer files in later stages.

## Basic pattern matching

`grep` should search for a match within a file, if a match is found, `grep` should print the line to stdout. If no match is found, `grep` should print nothing to stdout and exit with status code 1.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single file. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines and 1 if not.

```
[setup] $ echo "2024-01-01 ERROR: Database connection failed" > app.log
[setup] $ echo "[DEBUG] 4 errors found" > debug.log
$ grep "ERROR" app.log
2024-01-01 ERROR: Database connection failed
$ grep -E "\d+ errors? found" debug.log
[DEBUG] 4 errors found
$ grep -E "^\d{4}-\d{2}-\d{2} ERROR:" app.log
2024-01-01 ERROR: Database connection failed
$ grep -E ".* EROR" app.log
$ echo $?
1
```

## Notes

- The file is guaranteed to exist and be of a single line
- Output should contain the full line that matches the pattern

# Stage 2: Multiple-line file search

In this stage, you'll add support for pattern matching on the contents of a single file, which will consist of multiple lines.

## Basic pattern matching

`grep` should search for matches within a file, if a match is found, `grep` should print the line to stdout. `grep` should process the file line by line, and not error out on the first line that doesn't match the pattern. If no match is found in the entire file, `grep` should print nothing to stdout and exit with status code 1.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single file. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines and 1 if not.

```
[setup] $ rm app.log
[setup] $ echo "2024-01-01 ERROR: Database connection failed" > app.log
[setup] $ echo "2024-01-01 DEBUG: Query executed" >> app.log
[setup] $ echo "2024-01-01 ERROR: SQL syntax error" >> app.log
$ grep "DEBUG" app.log
2024-01-01 DEBUG: Query executed
$ grep -E "^\d{4}-\d{2}-\d{2} DEBUG:" app.log
2024-01-01 DEBUG: Query executed
$ grep -E ".* ERROR: .*" app.log
2024-01-01 ERROR: Database connection failed
2024-01-01 ERROR: SQL syntax error
$ grep -E ".* DEBUG: .* error" app.log
$ echo $?
1
```

## Notes

- The file is guaranteed to exist and be of multiple lines
- Output should contain the full lines that match the pattern

# Stage 3: Multiple-file search

In this stage, you'll add support for pattern matching on the contents of multiple files.

## Multi-file search

`grep` processes each file independently and handles results on a per-file basis.

The behavior follows these rules:

**File processing**: Files with matches will output all matching lines in their entirety to stdout with a `<filename>:` prefix. Files without matches produce no output but do not affect the exit code if other files contain matches. The filename used in the prefix included the path as passed to `grep`.

**Exit code behavior**: The exit code is determined by the overall operation result. Exit code 0 indicates at least one file contained matches. Exit code 1 indicates no matches were found in any existing file.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches across multiple files. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines and 1 if not.

```
[setup] $ echo "#include <stdio.h>" > main.c
[setup] $ echo "int main() {" >> main.c
[setup] $ echo "    printf(\"Hello World!\");" >> main.c
[setup] $ echo "    return 0;" >> main.c
[setup] $ echo "}" >> main.c
[setup] $ echo "#include <iostream>" > main.cpp
[setup] $ echo "using namespace std;" >> main.cpp
[setup] $ echo "int main() {" >> main.cpp
[setup] $ echo "    cout << \"C++ Program\" << endl;" >> main.cpp
[setup] $ echo "    return 0;" >> main.cpp
[setup] $ echo "}" >> main.cpp
[setup] $ echo "def main():" > script.py
[setup] $ echo "    database_host = \"localhost\"" >> script.py
[setup] $ echo "    database_user = \"admin\"" >> script.py
[setup] $ echo "    database_password = \"secret123\"" >> script.py
[setup] $ echo "" >> script.py
[setup] $ echo "if __name__ == \"__main__\":" >> script.py
[setup] $ echo "    main()" >> script.py
$ grep "include" main.c main.cpp
main.c:#include <stdio.h>
main.cpp:#include <iostream>
$ grep -E ".* main[\(\)]+" main.c main.c main.cpp script.py
main.c:int main() {
main.c:int main() {
main.cpp:int main() {
script.py:def main():
script.py:    main()
$ grep -E "class.*"  main.c main.cpp script.py
$ echo $?
1
```

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