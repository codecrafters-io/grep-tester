.PHONY: release build test test_with_bash copy_course_file

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))

release:
	git tag v$(next_version_number)
	git push origin main v$(next_version_number)

record_fixtures:
	CODECRAFTERS_RECORD_FIXTURES=true make test

build:
	go build -o dist/main.out ./cmd/tester

test:
	TESTER_DIR=$(shell pwd) go test -v ./internal/

test_and_watch:
	onchange '**/*' -- go test -v ./internal/

test_base_with_grep: build
	CODECRAFTERS_REPOSITORY_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\":\"cq2\",\"tester_log_prefix\":\"stage-1\",\"title\":\"Stage #1: Match a literal character\"},{\"slug\":\"oq2\",\"tester_log_prefix\":\"stage-2\",\"title\":\"Stage #2: Match digits\"},{\"slug\":\"mr9\",\"tester_log_prefix\":\"stage-3\",\"title\":\"Stage #3: Match alphanumeric characters\"},{\"slug\":\"tl6\",\"tester_log_prefix\":\"stage-4\",\"title\":\"Stage #4: Positive Character Groups\"},{\"slug\":\"rk3\",\"tester_log_prefix\":\"stage-5\",\"title\":\"Stage #5: Negative Character Groups\"},{\"slug\":\"sh9\",\"tester_log_prefix\":\"stage-6\",\"title\":\"Stage #6: Combining Character Classes\"},{\"slug\":\"rr8\",\"tester_log_prefix\":\"stage-7\",\"title\":\"Stage #7: Start of string anchor\"},{\"slug\":\"ao7\",\"tester_log_prefix\":\"stage-8\",\"title\":\"Stage #8: End of string anchor\"},{\"slug\":\"fz7\",\"tester_log_prefix\":\"stage-9\",\"title\":\"Stage #9: Match one or more times\"},{\"slug\":\"ny8\",\"tester_log_prefix\":\"stage-10\",\"title\":\"Stage #10: Match zero or one times\"},{\"slug\":\"zb3\",\"tester_log_prefix\":\"stage-11\",\"title\":\"Stage #11: Wildcard\"},{\"slug\":\"zm7\",\"tester_log_prefix\":\"stage-12\",\"title\":\"Stage #12: Alternation\"}]" \
	dist/main.out

test_backreferences_with_grep: build
	CODECRAFTERS_REPOSITORY_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\":\"sb5\",\"tester_log_prefix\":\"stage-13\",\"title\":\"Stage #13: Single Backreference\"},{\"slug\":\"tg1\",\"tester_log_prefix\":\"stage-14\",\"title\":\"Stage #14: Multiple Backreferences\"},{\"slug\":\"xe5\",\"tester_log_prefix\":\"stage-15\",\"title\":\"Stage #15: Nested Backreferences\"}]" \
	dist/main.out

test_file_search_with_grep: build
	CODECRAFTERS_REPOSITORY_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\":\"dr5\",\"tester_log_prefix\":\"stage-16\",\"title\":\"Stage #16: Single Line File Search\"},{\"slug\":\"ol9\",\"tester_log_prefix\":\"stage-17\",\"title\":\"Stage #17: Multi Line File Search\"},{\"slug\":\"is6\",\"tester_log_prefix\":\"stage-18\",\"title\":\"Stage #18: Multiple Files Search\"},{\"slug\":\"yx6\",\"tester_log_prefix\":\"stage-19\",\"title\":\"Stage #19: Recursive File Search\"}]" \
	dist/main.out

test_quantifiers_with_grep: build
	CODECRAFTERS_REPOSITORY_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\":\"ai9\",\"tester_log_prefix\":\"stage-20\",\"title\":\"Stage #20: Zero or more times\"},{\"slug\":\"wy9\",\"tester_log_prefix\":\"stage-21\",\"title\":\"Stage #21: Exact repetition\"},{\"slug\":\"hk3\",\"tester_log_prefix\":\"stage-22\",\"title\":\"Stage #22: Minimum repetition\"},{\"slug\":\"ug0\",\"tester_log_prefix\":\"stage-23\",\"title\":\"Stage #23: Range Repetition\"}]" \
	dist/main.out

test_printing_matches_with_grep: build
	CODECRAFTERS_REPOSITORY_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\":\"ku5\",\"tester_log_prefix\":\"stage-30\",\"title\":\"Stage #30: Print a single matching line\"},{\"slug\":\"pz6\",\"tester_log_prefix\":\"stage-31\",\"title\":\"Stage #31: Print multiple matching lines\"}]" \
	dist/main.out

test_multiple_matches_with_grep: build
	CODECRAFTERS_REPOSITORY_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\":\"cj0\",\"tester_log_prefix\":\"stage-40\",\"title\":\"Stage #40: Print single match\"},{\"slug\":\"ss2\",\"tester_log_prefix\":\"stage-41\",\"title\":\"Stage #41: Print multiple matches\"}, {\"slug\":\"bo4\",\"tester_log_prefix\":\"stage-42\",\"title\":\"Stage #42: Print multiple input lines\"}]" \
	dist/main.out

test_all: build
	make test_base_with_grep || true
	make test_backreferences_with_grep || true
	make test_file_search_with_grep || true
	make test_quantifiers_with_grep || true
	make test_printing_matches_with_grep || true
	make test_multiple_matches_with_grep || true

copy_course_file:
	hub api \
		repos/codecrafters-io/core/contents/data/courses/grep.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils
