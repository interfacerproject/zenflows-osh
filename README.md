<!--
Written and maintained by srfsh <info@dyne.org>.
Copyright (C) 2023 Dyne.org foundation <foundation@dyne.org>.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# zenflows-osh

A tool to clone and analyze git(1) repositories.


## Bulding and executing

You may either choose to build from source, build and use the Docker
image, or just use the pre-built Docker image on
https://ghcr.io/interfacerproject/zenflows-osh.

Building from source with Go might be a bit work due to the
dependencies.  Building from source with Docker is recommended as
the produced image will be self-contained.


### Building from source with Go

Bulding with Go from the source requires Go version 1.19 or later.
If you have the Go toolchain and a POSIX-compliant make(1)
implementation installed (GNU make(1) works), you can just run:

	make serve

which builds the development version of the service as an executable
named `zosh`.  You may use:

	make build.rel

to build the release version, and use the binary `zosh`.  Using
`make serve` again would build the development version, so use the
produced executable directly with `./zosh`.  At the moment, there's
no difference between the development and release versions.

If you choose to use this approach, you must have git(1),
[OSH tool](https://github.com/hoijui/osh-tool), and
[projvar](https://github.com/hoijui/projvar) (which is a dependency
of OSH) programs included in your $PATH.


### Building from source with Docker

Building with Docker requires nothing but the Docker tooling.  The
image produced will have all the dependencies needed to run this
service.

To build the image, you can run:

	docker build -t zenflows-osh:latest .

which will name the image "zenflosh-osh".  Then, you can run:

	docker run --rm -p PORT:7000 zenflows-osh

to start the service on port `PORT`.


### Using the pre-built image

You may choose to just use the pre-built image, which is found at
https://ghcr.io/interfacerproject/zenflows-osh.

Then, you can run:

	docker run --rm -p PORT:7000 ghcr.io/interfacerproject/zenflows-osh

to start the service on port `PORT`.

You may optionally use a docker-compose.yml template like this as well:

```
version: "3.8"
services:
  zosh:
    container_name: zosh
    image: ghcr.io/interfacerproject/zenflows-osh
    ports:
      # The service will be listening on port 3000 of the host
      # machine.
      - 3000:7000
    stdin_open: true
    tty: true
```

## Configuration

These only one configuration option at this moment, and it is `ADDR`.

`ADDR` is the address which the service should bind on.  It is of
the form `host:port`, where `host` can be any valid hostname or IP
addresses of both families, and `port` is a port number between 0
and 65535.  The `host` part could be omitted, which defaults to
listening on all IP addresses available on the host of both families

It has a default value of `:7000` when it is not provided.


## Usage

The usage schema is really easy.  The main idea is that the URL
defines a procedure, such as "/analyze", and each procedure takes
arguments in the form of a json, which is sent over the body of a
POST request.

If the procedure succeeds, a json of the form `{"ok": RESULT}`,
where `RESULT` is another json object or an array, will be returned
with 200 status code.

If the procedure fails, a json of the form `{"err": REASON}`,
where `REASON` is a string describing the error, will be returned with 500 status code.

If the procedure is not implemented (not found), a json of the form
`{"err": REASON}`, where `REASON` is a string describing that this
procedure doesn't exist, will be returned with 404 status code.

There's no other valid status code or response this service returns.

Here is an example:

	curl -XPOST -d'{"repo": "https://github.com/interfacerproject/zenflows-osh"}' https://zenflows-osh.interfacer.dyne.org/analyze

This will return `{"ok": RESULT}`, where `RESULT` is a json generated
by the OSH-tool.


## Procedures implemented

Here is a list of procedures you can use with this service.

### Analyze

The analyze procedure is accessed with `/analyze`.  it takes these parameters:

* `repo` - a URL to a repository

and returns the json data generated by the OSH-tool.
