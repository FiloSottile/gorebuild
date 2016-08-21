`gorebuild` will use DWARF symbols to figure out the import path a Go binary was built with, and will forcefully re-install it.

This is useful for example when changing GOROOT (for example updating Go with Homebrew), since the default GOROOT is embedded in the binary and it's critical for some static analysis tools.

Installation: `go get -u github.com/FiloSottile/gorebuild`

Usage: `gorebuild [-n] [binary ...]`

If invoked with `-n` it will only print the import paths.

If invoked without any arguments, it runs on all files in `$GOPATH/bin`.
