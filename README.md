## Overview

`ethrlp` is a Go library that implements encoding and decoding logic for Ethereum's [Recursive Length Prefix (RLP)
standard](https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/). RLP is a serialization method used
in Ethereum to encode arbitrarily nested arrays of binary data. This
library allows developers to encode various data types to RLP format and decode RLP-encoded data back into its original
form.

This library intentionally does not expose an `any` API for encoding, since it wants to avoid using reflection.
The consequence of this is that the package caller will need to manually encode specific struct fields, array values,
using the provided encode methods.

## Installation

You can install `ethrlp` using `go get`:

```shell
go get -u github.com/madz-lab/ethrlp
```

## Usage

For documentation on how to use the library, please reference the
adequate [Go Doc](https://pkg.go.dev/github.com/madz-lab/ethrlp) page.

## Benchmarks

```shell
goos: darwin
goarch: arm64
pkg: github.com/madz-lab/ethrlp
BenchmarkDecode_String_Short-14                 55288180                18.69 ns/op           24 B/op          1 allocs/op
BenchmarkDecode_String_Medium-14                62132836                19.18 ns/op           24 B/op          1 allocs/op
BenchmarkDecode_String_Long-14                  61043461                19.45 ns/op           24 B/op          1 allocs/op
BenchmarkDecode_Array_Small-14                  15084543                79.25 ns/op          136 B/op          4 allocs/op
BenchmarkDecode_Array_Medium-14                  2826644               431.5 ns/op           832 B/op         19 allocs/op
BenchmarkDecode_Array_Large-14                    983508              1243 ns/op            3216 B/op         51 allocs/op
BenchmarkDecode_Array_Nested_Short-14            8398964               142.8 ns/op           248 B/op          7 allocs/op
BenchmarkDecode_Array_Nested_Long-14             2572453               464.9 ns/op           728 B/op         22 allocs/op
BenchmarkEncode_String_Short-14                 76779115                15.06 ns/op           64 B/op          1 allocs/op
BenchmarkEncode_String_Medium-14                61054846                19.40 ns/op           64 B/op          1 allocs/op
BenchmarkEncode_String_Long-14                   9473983               116.8 ns/op          1152 B/op          1 allocs/op
BenchmarkEncode_Array_Small-14                  23130206                50.32 ns/op           69 B/op          4 allocs/op
BenchmarkEncode_Array_Medium-14                  4371036               275.7 ns/op           612 B/op         19 allocs/op
BenchmarkEncode_Array_Large-14                   1774267               670.8 ns/op          1668 B/op         49 allocs/op
BenchmarkEncode_Array_Nested_Short-14            9627250               125.0 ns/op           144 B/op          9 allocs/op
BenchmarkEncode_Array_Nested_Long-14             3124429               380.8 ns/op           720 B/op         25 allocs/op
PASS
ok      github.com/madz-lab/ethrlp      22.270s
```