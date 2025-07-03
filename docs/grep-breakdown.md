Grep / File Search

# Stage 1: Search a single-line file

In this stage, you'll add support for searching the contents of a file with a single line.

## File Search

When `grep` is given a file as an argument, it searches through the lines in the file and prints out matching lines. Example usage:

```bash
# This prints any lines that match search_pattern
$ grep -E "search_pattern" any_file.txt
This is a line that matches search_pattern
```

Matching lines are printed to stdout.

If any matching lines were found, grep exits with status code 0 (i.e. "success"). If no matching lines were found, grep exits with status code 1.

In this stage, we'll test searching through a file that contains a single line. We'll get to handling multi-line files in later stages.

## Tests

The tester will create some test files and then execute multiple commands to find matches in those files. For example:

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

## Multiple matches within a file

When searching through a multi-line file, `grep` processes each line individually and prints all matching lines to stdout, each on its own line. Example usage:

```bash
# This prints any lines that match search_pattern
$ grep -E "search_pattern" multi_line_file.txt
This is a line that matches search_pattern
This is another line that matches search_pattern
```

All matching lines are printed to stdout.

If any matching lines were found, grep exits with status code 0 (i.e. "success"). If no matching lines were found, grep exits with status code 1.

## Tests

The tester will create some test files and then execute multiple commands to find matches in those files. For example:

```bash
# Create test file
$ echo "cherry" > fruits.txt
$ echo "banana" >> fruits.txt
$ echo "grape" >> fruits.txt
$ echo "blueberry" >> fruits.txt

# This must print the matched lines to stdout and exit with code 0
$ ./your_program.sh -E ".*erry" fruits.txt
cherry
blueberry

# This must print nothing to stdout and exit with code 1
$ ./your_program.sh -E "app.*" fruits.txt
```

The tester will verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there aren't.

## Notes

- The file is guaranteed to exist and contain multiple lines
- Output should contain the full lines that match the pattern

# Stage 3: Search multiple files

In this stage, you'll add support for searching the contents of multiple files.

## Searching multiple files

When searching multiple files, `grep` processes each file individually and prints all matching lines to stdout with a `<filename>:` prefix. Example usage:

```bash
# This prints any lines that match search_pattern from multiple files
$ grep -E "search_pattern" file1.txt file2.txt
file1.txt:This is a line that matches search_pattern
file2.txt:Another line that matches search_pattern
```

Matching lines are printed to stdout with filename prefixes.

If any matching lines were found, grep exits with status code 0 (i.e. "success"). If no matching lines were found, grep exits with status code 1.

## Tests

The tester will create some test files and then execute multiple commands to find matches in those files. For example:

```bash
# Create test files
$ echo "cherry" > fruits.txt
$ echo "blueberry" >> fruits.txt
$ echo "celery" > vegetables.txt
$ echo "carrot" >> vegetables.txt

# This must print the matched lines to stdout and exit with code 0
$ ./your_program.sh -E ".*ry$" fruits.txt vegetables.txt
fruits.txt:cherry
fruits.txt:blueberry
vegetables.txt:celery

# This must print nothing to stdout and exit with code 1
$ ./your_program.sh -E "on.*on" fruits.txt vegetables.txt
```

The tester will verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there aren't.

# Stage 4: Recursive search

In this stage, you'll add support for searching through files in a given directory and its subdirectories recursively with the `-r` flag.

## Recursive search

When `grep` is passed the `-r` flag, it searches through the given directory and its subdirectories recursively. It processes each file line by line and prints all matching lines to stdout with a `<filename>:` prefix. Example usage:

```bash
# This prints any lines that match search_pattern from multiple files
$ grep -r -E "search_pattern" directory/
directory/file1.txt:This is a line that matches search_pattern
directory/file2.txt:Another line that matches search_pattern
```

Matching lines are printed to stdout with filename prefixes.

If any matching lines were found, grep exits with status code 0 (i.e. "success"). If no matching lines were found, grep exits with status code 1.

## Tests

The tester will create some test files and then multiple commands to find matches in those files. For example:

```bash
# Create test files
$ mkdir -p dir/subdir
$ echo "pear" > dir/fruits.txt
$ echo "strawberry" >> dir/fruits.txt
$ echo "celery" > dir/subdir/vegetables.txt
$ echo "carrot" >> dir/subdir/vegetables.txt
$ echo "cucumber" > dir/vegetables.txt
$ echo "corn" >> dir/vegetables.txt

# This must print the matched lines to stdout and exit with code 0
$ ./your_program.sh -r -E ".*er" dir/
dir/fruits.txt:strawberry
dir/subdir/vegetables.txt:celery
dir/vegetables.txt:cucumber

# This must print the matched lines to stdout and exit with code 0
$ ./your_program.sh -r -E ".*ar" dir/
dir/fruits.txt:pear
dir/subdir/vegetables.txt:carrot

# This must print nothing to stdout and exit with code 1
$ ./your_program.sh -r -E "orange" dir/
```

The tester will verify that all matching lines are printed to stdout. It'll also verify that the exit code is 0 if there are matching lines, and 1 if there aren't.

## Notes

- GNU grep doesn't guarantee the sorting order of output; it processes files in filesystem order. Your `grep` can output matching lines in any order.
- The filepath prefix is relative to the directory passed as an argument to the `-r` flag
---
