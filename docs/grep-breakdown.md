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
$ grep -E ".* main[\(\)]+" main.c ./main.c main.cpp script.py
main.c:int main() {
./main.c:int main() {
main.cpp:int main() {
script.py:def main():
script.py:    main()
$ grep -E "class.*"  main.c main.cpp script.py
$ echo $?
1
```

# Stage 4: Single-directory recursive search

In this stage, you'll add support for searching through files in a given directory and its subdirectories recursively with the `-r` flag.

## Recursive search

The `-r` flag enables recursive searching through directories and their subdirectories. `grep` should search for matches in each file it finds, and process the file line by line. Each matching line should be prefixed with the relative path to the file `<filename>:` (the filepath is relative from the directory passed to `grep` as input).

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches in a single directory. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines and 1 if not.

```
[setup] $ rm -rf logs/
[setup] $ mkdir -p logs/deeply/nested
[setup] $ echo "ERROR: Database connection failed" > logs/app.log
[setup] $ echo "ERROR: Nested error" > logs/deeply/file.log
[setup] $ echo "WARN: Might be a warning" >> logs/deeply/file.log
[setup] $ echo "INFO: This is alright!" >> logs/deeply/file.log
[setup] $ echo "2024-01-01 ERROR: Database connection failed" > logs/deeply/nested/app.log
[setup] $ echo "2024-01-01 INFO: Server started successfully" >> logs/deeply/nested/app.log
[setup] $ echo "2024-01-01 DEBUG: Processing user request" >> logs/deeply/nested/app.log
[setup] $ echo "2024-01-01 ERROR: SQL syntax error in query" >> logs/deeply/nested/app.log
$ grep -r "ERROR" logs/
logs/deeply/file.log:ERROR: Nested error
logs/deeply/nested/app.log:2024-01-01 ERROR: Database connection failed
logs/deeply/nested/app.log:2024-01-01 ERROR: SQL syntax error in query
logs/app.log:ERROR: Database connection failed
$ cd logs
$ grep -r -E "^\d{4}-\d{2}-\d{2} ERROR:" .
logs/deeply/nested/app.log:2024-01-01 ERROR: Database connection failed
logs/deeply/nested/app.log:2024-01-01 ERROR: SQL syntax error in query
$ cd logs
$ grep -r -E "Database.*connection.*failed?"
logs/deeply/nested/app.log:2024-01-01 ERROR: Database connection failed
logs/app.log:ERROR: Database connection failed
$ cd ..
$ grep -r -E "(success|info)$" .
$ echo $?
1
```

## Notes

- `-r` doesn't follow recursive symlinks (we won't test for symlinks at all)
- GNU Grep doesn't guarantee the sorting order of the output, it processes the files in the order the underlying filesystem returns them. You can return the output in chronological order if you want. We won't test for this.
- If no directory is provided with `-r`, `grep` runs the search in the current working directory.

# Stage 5: Multiple-directory recursive search

In this stage, you'll add support for searching through files in multiple directories and their subdirectories recursively with the `-r` flag.

## Multiple-directory recursive search

The `-r` flag enables recursive searching through multiple directories and their subdirectories. `grep` should search for matches in each directory and file it finds across all specified directories, processing each file line by line. Each matching line should be prefixed with the relative path to the file `<filename>:` (the filepath is relative from each directory passed to `grep` as input). `grep` handles each directory independently, and the output is not sorted.

## Tests

The tester will execute your program like this:

```bash
./your_program.sh
```

It will then run multiple `grep` commands to find matches across multiple directories. The tester will then verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines and 1 if not.

```
[setup] $ rm -rf logs/ src/
[setup] $ mkdir -p logs/app src/utils
[setup] $ echo "ERROR: Database connection failed" > logs/app.log
[setup] $ echo "INFO: Server started" >> logs/app.log
[setup] $ echo "ERROR: Authentication failed" > logs/app/auth.log
[setup] $ echo "function validateUser(id)" > src/auth.js
[setup] $ echo "function processData(input)" > src/utils/helper.js
[setup] $ echo "ERROR: Processing failed" >> src/utils/helper.js
$ grep -r "ERROR" logs/ src/
logs/app/auth.log:ERROR: Authentication failed
logs/app.log:ERROR: Database connection failed
src/utils/helper.js:ERROR: Processing failed
$ grep -r -E "function\s+\w+" logs/ src/
src/auth.js:function validateUser(id)
src/utils/helper.js:function processData(input)
$ grep -r "nonexistent" logs/ src/
$ echo $?
1
$ cd logs
$ grep -r -E "ERROR: \w+ failed" . app ../logs/app/auth.log
./app/auth.log:ERROR: Authentication failed
./app.log:ERROR: Database connection failed
app/auth.log:ERROR: Authentication failed
../logs/app/auth.log:ERROR: Authentication failed
```

## Notes

- Each directory maintains its own relative path context

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