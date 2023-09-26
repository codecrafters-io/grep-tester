.PHONY: release build test test_with_bash copy_course_file

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))

release:
	git tag v$(next_version_number)
	git push origin main v$(next_version_number)

build:
	go build -o dist/main.out ./cmd/tester

test:
	TESTER_DIR=$(shell pwd) go test -v ./internal/

test_and_watch:
	onchange '**/*' -- go test -v ./internal/

test_with_grep: build
	CODECRAFTERS_SUBMISSION_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{"slug":"init","tester_log_prefix":"stage-1","title":"Stage #1: Match a literal character"},{"slug":"match_digit","tester_log_prefix":"stage-2","title":"Stage #2: Match digits"},{"slug":"match_alphanumeric","tester_log_prefix":"stage-3","title":"Stage #3: Match alphanumeric characters"},{"slug":"positive_character_groups","tester_log_prefix":"stage-4","title":"Stage #4: Positive Character Groups"},{"slug":"negative_character_groups","tester_log_prefix":"stage-5","title":"Stage #5: Negative Character Groups"},{"slug":"combining_character_classes","tester_log_prefix":"stage-6","title":"Stage #6: Combining Character Classes"},{"slug":"start_of_string_anchor","tester_log_prefix":"stage-7","title":"Stage #7: Start of string anchor"},{"slug":"end_of_string_anchor","tester_log_prefix":"stage-8","title":"Stage #8: End of string anchor"},{"slug":"one_or_more_quantifier","tester_log_prefix":"stage-9","title":"Stage #9: Match one or more times"},{"slug":"zero_or_one_quantifier","tester_log_prefix":"stage-10","title":"Stage #10: Match zero or one times"},{"slug":"wildcard","tester_log_prefix":"stage-11","title":"Stage #11: Wildcard"},{"slug":"alternation","tester_log_prefix":"stage-12","title":"Stage #12: Alternation"}]" \
	dist/main.out

copy_course_file:
	hub api \
		repos/codecrafters-io/core/contents/data/courses/grep.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils
