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

test_with_grep: build
	CODECRAFTERS_REPOSITORY_DIR=$(shell pwd)/internal/test_helpers/pass_all \
	CODECRAFTERS_TEST_CASES_JSON="[{\"slug\":\"cq2\",\"tester_log_prefix\":\"stage-1\",\"title\":\"Stage #1: Match a literal character\"},{\"slug\":\"oq2\",\"tester_log_prefix\":\"stage-2\",\"title\":\"Stage #2: Match digits\"},{\"slug\":\"mr9\",\"tester_log_prefix\":\"stage-3\",\"title\":\"Stage #3: Match alphanumeric characters\"},{\"slug\":\"tl6\",\"tester_log_prefix\":\"stage-4\",\"title\":\"Stage #4: Positive Character Groups\"},{\"slug\":\"rk3\",\"tester_log_prefix\":\"stage-5\",\"title\":\"Stage #5: Negative Character Groups\"},{\"slug\":\"sh9\",\"tester_log_prefix\":\"stage-6\",\"title\":\"Stage #6: Combining Character Classes\"},{\"slug\":\"rr8\",\"tester_log_prefix\":\"stage-7\",\"title\":\"Stage #7: Start of string anchor\"},{\"slug\":\"ao7\",\"tester_log_prefix\":\"stage-8\",\"title\":\"Stage #8: End of string anchor\"},{\"slug\":\"fz7\",\"tester_log_prefix\":\"stage-9\",\"title\":\"Stage #9: Match one or more times\"},{\"slug\":\"ny8\",\"tester_log_prefix\":\"stage-10\",\"title\":\"Stage #10: Match zero or one times\"},{\"slug\":\"zb3\",\"tester_log_prefix\":\"stage-11\",\"title\":\"Stage #11: Wildcard\"},{\"slug\":\"zm7\",\"tester_log_prefix\":\"stage-12\",\"title\":\"Stage #12: Alternation\"}]" \
	dist/main.out

copy_course_file:
	hub api \
		repos/codecrafters-io/core/contents/data/courses/grep.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils

setup:
	curl -LO https://github.com/BurntSushi/ripgrep/releases/download/14.0.2/ripgrep_14.0.2-1_amd64.deb
	sudo dpkg -i ripgrep_14.0.2-1_amd64.deb
	rm ripgrep_14.0.2-1_amd64.deb

setup_bsdgrep:
	sudo apt update
	sudo apt install -y build-essential curl git
	git clone https://github.com/arp242/bsdgrep.git
	pwd && ls -laH
	cd bsdgrep
	pwd && ls -laH
	ls -laH bsdgrep/
	cd bsdgrep && pwd && ls -la
	./update.sh
	sed -i 's/#error.*getprogname.*/return \"grep\";/' progname.c
	sed -i 's/warnc(/warn(/g' util.c
	sed -i 's/warn(p->fts_errno,/warn(/g' util.c
	rm freebsd.c
	make
	make install
	which grep
	grep --version
	mv /usr/bin/grep /usr/bin/grep.gnu && ln -sf /usr/local/bin/grep /usr/bin/grep
	grep --version