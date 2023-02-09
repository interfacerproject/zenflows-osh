# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2023 Dyne.org foundation <foundation@dyne.org>.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

ARG GOVER=1.19
ARG ALPVER=3.17


FROM golang:$GOVER-alpine$ALPVER AS golang
WORKDIR /app
RUN apk --no-cache add make
COPY . .
RUN make build.rel


FROM alpine:$ALPVER AS osh-tool
WORKDIR /app/tmp
RUN apk --no-cache add git nim nimble gcc musl-dev pcre-dev openssl1.1-compat-libs-static
RUN git clone --recurse-submodules https://github.com/hoijui/osh-tool .
RUN nimble build -y \
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
WORKDIR /app/tmp
RUN wget -qO projvar.tgz \
	https://github.com/hoijui/projvar/releases/download/0.16.0/projvar-0.16.0-x86_64-unknown-linux-musl.tar.gz
RUN tar -xzf projvar.tgz && mv projvar-*-x86_64-unknown-linux-musl/projvar ..


FROM alpine:$ALPVER
WORKDIR /app
ARG PORT=7000
ENV ADDR=:$PORT
EXPOSE $PORT
ARG USER=zosh
ARG GROUP=$USER
ENV PATH="$PATH:/app/bin"

RUN apk --no-cache add git

RUN addgroup -S "$GROUP" && adduser -SG"$GROUP" "$USER"

COPY --from=projvar /app/projvar ./bin/
COPY --from=osh-tool /app/osh ./bin/
COPY --from=golang --chown="$USER:$GROUP" /app/zosh ./

CMD ./zosh
