# Written and maintained by srfsh <info@dyne.org>.
# Copyright (C) 2023 Dyne.org foundation <foundation@dyne.org>.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

.POSIX:
.SUFFIXES:

# in order for it to not conflict with osh-tool binary
PROG=zosh

all: build

build: build.dev

build.dev:
	go build -o ${PROG}

build.rel:
	go build -o ${PROG}

serve: build.dev
	./${PROG}

test: build.dev
	go test ./...

clean:
	rm -rf ${PROG}

help:
	@echo "build|build.dev:	build for development (the default target)"
	@echo "build.rel:	build for release"
	@echo "serve:		build and execute the program"
	@echo "test:		run the tests against the development binary"
	@echo "clean:		clean any generated file"
	@echo "help:		print this text"

.PHONY: all build build.dev build.rel test clean help
