**Deprecated**: `gorebuild` does not support Go modules, and the need to update the built-in default GOROOT has mostly gone away, in favor of using [`go/packages`](https://pkg.go.dev/golang.org/x/tools/go/packages) which invokes the go tool from `$PATH`.

`gorebuild` uses symbol tables to figure out the import path of a Go binary, re-installs it.

This is useful for example when changing GOROOT (for example updating Go with Homebrew), since the default GOROOT is embedded in the binary and it's critical for some static analysis tools.

Installation: `go get -u github.com/FiloSottile/gorebuild`

Usage: `gorebuild [-n] [binary ...]`

If invoked with `-n` it will only print the import paths.

If invoked without any arguments, it runs on all files in `$GOPATH/bin`.
