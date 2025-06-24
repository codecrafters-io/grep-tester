# Grep.

## File Search.
Ref: https://www.gnu.org/software/grep/manual/grep.html

1. Match on single file contents.
```
> grep "ERROR" main.c
ERROR: Buffer overflow detected
> Pattern search.
> Regex search.
```

2. Match across multiple files.
```
> grep "include" main.c program.cpp
main.c:#include <stdio.h>
program.cpp:#include <iostream>
> Pattern search.
> Regex search.
```

3. Match against directory contents without -r or -d.
```
> grep "include" src/
grep: src/: Is a directory
```

4. Match against directory contents recursively.
Note: -r doesn't follow recursive symlinks. -R does. We use -r.
```
> grep -r "ERROR" logs
logs/large.log:ERROR_LINE
logs/app.log:ERROR: Database connection failed
```

## On the fence about these.

5. --include for inclusion patterns.
```
> grep -r --include="*.c" "include" src
src/c/main.c:#include <stdio.h>
```

6. --exclude for exclusion patterns.
```
> grep -r --exclude="*.log" "ERROR" docs
docs/readme.txt:Some lines have ERRORS.
```

7. Files with spaces, quotes, newlines.
```
> ggrep "PATTERN" docs/file\ with\ spaces.txt
PATTERN
> ggrep "PATTERN" docs/file$'\n'with$'\n'newlines.txt
PATTERN
```

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