# go-up

go-up is a small utility to install and update Go versions.

## Installing

Grab a precompiled release from the [Releases](releases) page, and place in a
directory available in your `$PATH`.

## Using

To list remote available versions, use `goup list-remote`:

```
$ goup list-remote
 - 1.21.0 (stable, compatible)
    Archs: aix-ppc64, darwin-amd64, darwin-arm64, dragonfly-amd64, freebsd-386, freebsd-amd64, freebsd-arm64, freebsd-armv6l, freebsd-riscv64, illumos-amd64, linux-386, linux-amd64, linux-arm64, linux-armv6l, linux-loong64, linux-mips, linux-mips64, linux-mips64le, linux-mipsle, linux-ppc64, linux-ppc64le, linux-riscv64, linux-s390x, netbsd-386, netbsd-amd64, netbsd-arm64, netbsd-armv6l, openbsd-386, openbsd-amd64, openbsd-arm64, openbsd-armv6l, plan9-386, plan9-amd64, plan9-armv6l, solaris-amd64, windows-386, windows-amd64, windows-arm64, windows-armv6l
 - 1.20.7 (stable, compatible)
    Archs: darwin-amd64, darwin-arm64, freebsd-386, freebsd-amd64, linux-386, linux-amd64, linux-arm64, linux-armv6l, linux-ppc64le, linux-s390x, windows-386, windows-amd64, windows-arm64
```

To install a version, use `goup install VERSION`. For instance, to install `1.21.0`:

```
$ goup install 1.21.0
```

You can also use `latest` to install the latest version:

```
$ goup install latest
Querying go.dev... OK.
About to download and install go1.21.0. Continue? (y/N) y

Downloading version go1.21.0... OK.
Checking file integrity... OK.
Preparing to decompress archive...OK.
Decompressing SDK... OK.
Checking installation... OK. go version go1.21.0 darwin/amd64
Instalation completed. Use goup use 1.21.0 to activate the new version.
```

To list installed versions, use `goup list`:

```
$ goup list
 - 1.21.0 active healthy
```

To use a specific installed version, use `goup use`:

```
$ goup use 1.21.0
Activated go-1.21.0.
```

Finally, make sure that `$HOME/.go/current/bin` is present in your `$PATH`. It is
advised to place that directory in the beginning of the list.

## LICENSE

```
MIT License

Copyright (c) 2023 Victor Gama de Oliveira

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
