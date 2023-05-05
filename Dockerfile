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

# syntax=docker/dockerfile:1

ARG ALP_VER=3.17
ARG GO_VER=1.19

ARG OKH_TOOL_VER=0.4.4
ARG PROJVAR_VER=0.16.2
ARG MLE_VER=0.23.0
ARG OSH_DIR_STD_VER=0.7.0
ARG OSH_TOOL_VER=0.4.0

FROM golang:$GO_VER-alpine$ALP_VER AS golang
RUN apk --no-cache add make
WORKDIR /app
ENV GONOPROXY=
COPY . .
RUN make build.rel


FROM alpine:$ALP_VER as okh-tool
WORKDIR /app/tmp
ARG OKH_TOOL_VER
RUN wget -qO- -- \
		"https://github.com/OPEN-NEXT/LOSH-OKH-tool/releases/download/$OKH_TOOL_VER/okh-tool-$OKH_TOOL_VER-x86_64-unknown-linux-musl.tar.gz" \
	| tar -xzf -
RUN  mv -- "okh-tool-$OKH_TOOL_VER-x86_64-unknown-linux-musl/okh-tool" ..


FROM alpine:$ALP_VER AS projvar
WORKDIR /app/tmp
ARG PROJVAR_VER
RUN wget -qO- -- \
		"https://github.com/hoijui/projvar/releases/download/$PROJVAR_VER/projvar-$PROJVAR_VER-x86_64-unknown-linux-musl.tar.gz" \
	| tar -xzf -
RUN  mv -- "projvar-$PROJVAR_VER-x86_64-unknown-linux-musl/projvar" ..


FROM alpine:$ALP_VER AS mle
WORKDIR /app/tmp
ARG MLE_VER
RUN wget -qO- -- \
		"https://github.com/hoijui/mle/releases/download/$MLE_VER/mle-$MLE_VER-x86_64-unknown-linux-musl.tar.gz" \
	| tar -xzf -
RUN mv -- "mle-$MLE_VER-x86_64-unknown-linux-musl/mle" ..


FROM alpine:$ALP_VER AS osh-dir-std
WORKDIR /app/tmp
ARG OSH_DIR_STD_VER
RUN wget -qO- -- \
		"https://github.com/hoijui/osh-dir-std-rs/releases/download/$OSH_DIR_STD_VER/osh-dir-std-$OSH_DIR_STD_VER-x86_64-unknown-linux-musl.tar.gz" \
	| tar -xzf -
RUN mv -- "osh-dir-std-$OSH_DIR_STD_VER-x86_64-unknown-linux-musl/osh-dir-std" ..


FROM alpine:$ALP_VER AS osh-tool
RUN apk --no-cache add git nim nimble gcc musl-dev pcre-dev openssl1.1-compat-libs-static
WORKDIR /app/tmp
ARG OSH_TOOL_VER
RUN git clone -q --depth 1 --recurse-submodules -b "$OSH_TOOL_VER" \
	https://github.com/hoijui/osh-tool .
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
	--passL:-lcrypto
RUN mv build/osh ..


FROM alpine:$ALP_VER
RUN apk --no-cache add git reuse

ARG USER=zosh
ARG GROUP=$USER
RUN addgroup -S -- "$GROUP"
RUN adduser -SG "$GROUP" -- "$USER"

WORKDIR /app
COPY --from=okh-tool --chown=$USER:$GROUP /app/okh-tool ./bin/
COPY --from=projvar --chown=$USER:$GROUP /app/projvar ./bin/
COPY --from=mle --chown=$USER:$GROUP /app/mle ./bin/
COPY --from=osh-dir-std --chown=$USER:$GROUP /app/osh-dir-std ./bin/
COPY --from=osh-tool --chown=$USER:$GROUP /app/osh ./bin/
COPY --from=golang --chown=$USER:$GROUP /app/zosh ./

ENV PATH=$PATH:/app/bin

ARG PORT="7000"
EXPOSE "$PORT"
ENV ADDR=:$PORT

USER $USER:$GROUP

CMD ./zosh
