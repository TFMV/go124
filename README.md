# Go 1.24 Features Demo

This repository contains a demonstration of new features and improvements introduced in Go 1.24.

## Features Demonstrated

- Generic type aliases
- CGO improvements (noescape and nocallback annotations)
- Improved finalizers
- New crypto packages (HKDF, PBKDF2, SHA3)
- Directory-limited filesystem access
- Bytes and strings iterators
- New encoding interfaces (TextAppender and BinaryAppender)
- netip encoding interfaces
- Regexp TextAppender interface
- Runtime GOROOT deprecation notice
- Text template range over integer sequence
- math/big encoding TextAppender
- math/rand improvements
- sync.Map improvements
- log/slog DiscardHandler
- time encoding interfaces
- Experimental testing/synctest
- go/types iterator methods
- maphash comparable and WriteComparable

## Requirements

- Go 1.24 or later

## Running the Demo

To run the demo, simply execute:

```bash
go run go1.24_demo.go
```

## Notes

- Some features (like testing/synctest) require additional setup or environment variables
- CGO must be enabled to run the CGO-related demos
- Some features are best demonstrated in test files or with specific toolchain flags

## License

MIT License.
