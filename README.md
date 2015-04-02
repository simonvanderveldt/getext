# Getext
`getext` is a simple standalone binary written in Go that gets a `.tar.gz` file from a http(s) URL and extracts it to the current directory.
Primarily meant for use inside a (Docker) container.

Born out of annoyance with the fact that the `ADD` Dockerfile command doesn't extract remote resources and installing wget, tar and gzip inside your Dockerfile just to download and extract some files is way overkill.


## How to use
Usage is very straighforward
```
getext https://<domain>/filename.tar.gz
```

### Dockerfile examples
To work around needing wget, tar and gzip to download `getext` we make use of Docker's `ADD` Dockerfile command to add the `getext` binary to an image.<br>
`getext` binaries are provided using [GitHub Releases](https://github.com/simonvanderveldt/getext/releases).<br>
Note that ADDing `getext` in a Dockerfile isn't ideal either because every `ADD` command busts the cache :(

Currently the only supported Linux distro is [Alpine Linux](https://www.alpinelinux.org) using Gliderlabs [docker-alpine](https://github.com/gliderlabs/docker-alpine) Docker image.


__Dockerfile example__
```
FROM gliderlabs/alpine:3.1
ADD https://github.com/simonvanderveldt/getext/releases/download/v0.1.0/getext /usr/bin/
RUN ["chmod", "+x", "/usr/bin/getext"]
RUN ["getext", "https://<domain>/<filename>.tar.gz"]
```

Note that to get files from https sites you will need the necessary ca-certificates. There are two options to get those:

- Install using `apk add ca-certificates`, adds +- 9MB to the image
- Add `ADD https://github.com/bagder/ca-bundle/raw/master/ca-bundle.crt /etc/ssl/certs/ca-certificates.crt` to your Dockerfile, adds +- 250KB to the image
  - This makes use of cURL's awesome [converted CA certs from mozilla.org](http://curl.haxx.se/docs/caextract.html).

__Full Dockerfile example__
```
FROM gliderlabs/alpine:3.1
ADD https://github.com/simonvanderveldt/getext/releases/download/v0.1.0/getext /usr/bin/
RUN ["chmod", "+x", "/usr/bin/getext"]
ADD https://github.com/bagder/ca-bundle/raw/master/ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
RUN ["getext", "https://<domain>/<filename>.tar.gz"]
```
