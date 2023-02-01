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

ARG ALPVER=3.17

FROM golang:1.19-alpine$ALPVER AS golang
WORKDIR /app
RUN apk --no-cache add make
ADD . .
RUN make build.rel

FROM alpine:$ALPVER AS nim
WORKDIR /app
RUN apk --no-cache add git nim nimble gcc musl-dev pcre-dev openssl1.1-compat-libs-static
RUN git clone --recurse-submodules https://github.com/hoijui/osh-tool tmp
RUN cd tmp && nimble build -y \
	-d:release \
	--opt:speed \
	--passL:-static \
	--passL:-no-pie \
	-d:usePcreHeader \
	--passL:-lpcre \
	--dynlibOverride:ssl \
	--passL:-lssl \
	--dynlibOverride:crypto \
	--passL:-lcrypto && mv build/osh ..

FROM alpine:$ALPVER AS projvar
WORKDIR /app
RUN wget -qO projvar.tgz \
	https://github.com/hoijui/projvar/releases/download/0.16.0/projvar-0.16.0-x86_64-unknown-linux-musl.tar.gz
RUN mkdir tmp && cd tmp \
	&& tar -xzf ../projvar.tgz >/dev/null 2>&1 \
	&& mv projvar-*-x86_64-unknown-linux-musl/projvar ..

FROM alpine:$ALPVER
WORKDIR /app
EXPOSE 8000/tcp
ENV PATH="$PATH:/app/bin"
RUN apk --no-cache add git
COPY --from=golang /app/zosh .
COPY --from=nim /app/osh ./bin/
COPY --from=projvar /app/projvar ./bin/
CMD ["./zosh"]
