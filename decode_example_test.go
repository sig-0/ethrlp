package ethrlp

import "fmt"

func ExampleDecodeBytes() {
	// Byte value
	value, err := DecodeBytes([]byte{0x00})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%X %s\n", value.GetValue(), value.GetType())

	// Output:
	// 00 Bytes
}
