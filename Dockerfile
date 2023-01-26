ARG ALPVER=3.17

FROM golang:1.19-alpine$ALPVER AS golang
WORKDIR /app
ADD . .
RUN go build -o losh-import .

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
EXPOSE 80/tcp 443/tcp 443/udp
ENV PATH="$PATH:/app/bin"
RUN apk --no-cache add git
COPY --from=golang /app/losh-import .
COPY --from=nim /app/osh ./bin/
COPY --from=projvar /app/projvar ./bin/
CMD ["./losh-import"]
