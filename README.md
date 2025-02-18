# Go 1.24 Features Demo

This repository contains a demonstration of new features and improvements introduced in Go 1.24.

The code is intended to be self documenting.

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

## Output

The following is a sample output from the demo:

```bash
=== Go 1.24 Demo ===
Generic Type Alias (MySlice[int]): [1 2 3 4 5]
Derived key (HKDF): 283b1390f465a3f0b20cc90fee5c1ef1bf47fcd8393a9315547d1c63911c5b31
Derived key (PBKDF2): 6e52397ce1f677f36df6fd486bbbd5f611c264951ef1f4ebecaf1140a614d05e
SHA3-256 digest: 644bcc7e564373040999aac89e7622f3ca71fba1d972fd94a31c3bfbf24e3938
Files in limited FS:
 - example.txt
Iterating over lines (using bytes.Split):
line1
line2
line3
Iterating over fields (using strings.Fields):
foo
bar
baz
Encoding append result: demoStruct(123)
netip.Addr appended text: 192.0.2.1
Regexp appended text: a*b
Note: runtime.GOROOT is deprecated; use 'go env GOROOT' instead.
Template output: Numbers: 1 2 3 4 5 
big.Int appended text: 12345678901234567890
Random number (rand.New): 7382721707735833545
Iterating over sync.Map:
  key=key1, value=100
  key=key2, value=200
slog.DiscardHandler demo: In production, a DiscardHandler would discard logs.
Experimental synctest demo: See tests built with GOEXPERIMENT=synctest for usage.
go/types iterator demonstration: Use the Variables() method on tuples, etc.
Hash for key "myKey": 991567377541123890
=== Go 1.24 Demo End ===
```

## Note

Some features are not included in the demo.

## License

MIT License.
