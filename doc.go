// Package ethrlp is a Go library that implements encoding and decoding logic for Ethereum's Recursive Length Prefix
// (RLP) standard. RLP is a serialization method used in Ethereum to encode arbitrarily nested arrays of binary data.
// This library allows developers to encode various data types to RLP format and decode RLP-encoded data back into its
// original form.
//
// This library intentionally does not expose an `any` API for encoding, since it wants to avoid using reflection.
// The consequence of this is that the package caller will need to manually encode specific struct fields, array values,
// using the provided encode methods.
package ethrlp
