package ethrlp

// EncodeByte encodes a single byte to RLP
func EncodeByte(input byte) []byte {
	// If input is a single byte in the [0x00, 0x7f] range,
	// it itself is the RLP encoding
	if input <= 0x7f {
		return []byte{input}
	}

	// If input is a byte in the [0x80, 0xff] range,
	// RLP encoding will concatenate 0x81 with the byte
	return []byte{
		0x81,
		input,
	}
}

// EncodeBytes encodes a byte array to RLP
func EncodeBytes(input []byte) []byte {
	// If the input is a non-value
	// (uint(0), []byte{}, string(""), empty pointer...),
	// RLP encoding is 0x80
	if len(input) == 0 {
		return []byte{0x80}
	}

	// If input is more than 55 bytes long,
	// the RLP encoding consists of 3 parts:
	// - The first part is a single byte with value 0xb7
	// plus the length in bytes of the second part
	// - The second part is hex value of the length of the string
	// - The third part are the actual input bytes
	if len(input) > 55 {
		return encodeLongBytes(input)
	}

	// If input is between with 2–55 bytes long,
	// the RLP encoding consists of 2 parts:
	// - A single byte with value 0x80 plus the length of the second part,
	// - The second part are the actual input bytes
	return encodeShortBytes(input)
}

// encodeShortBytes encodes an input byte array that is <=55B long
func encodeShortBytes(input []byte) []byte {
	// The resulting RLP encoding is the entire input bytes array
	// along with the single byte denoting the length
	result := make([]byte, 0, 1+len(input))

	// The first byte is the sum of 0x80,
	// and the length of the input bytes
	result[0] = 0x80 + byte(len(input)) // length is <= 0xff

	// The rest of the RLP encoding
	// are the actual concatenated bytes
	return append(result, input...)
}

// encodeLongBytes encodes an input byte array that is >55B long
func encodeLongBytes(input []byte) []byte {
	// The resulting RLP encoding is the entire input bytes array
	// along with the single byte denoting the length,
	// and the actual length bytes
	lengthBytes := convertIntToHexArray(len(input))
	result := make([]byte, 0, 1+len(lengthBytes)+len(input))

	// The first byte is the sum of 0xb7
	// and the length of the incoming next bytes
	result[0] = 0xb7 + byte(len(lengthBytes))

	// The second part of the RLP encoding are the length bytes
	// of the input value
	result = append(result, lengthBytes...)

	// The rest of the RLP encoding
	// are the actual concatenated bytes
	return append(result, input...)
}

// convertIntToHexArray converts the integer value
// to an array of hex values (changes representation)
func convertIntToHexArray(length int) []byte {
	// Allocate a byte slice with capacity for 8 bytes (64 bits)
	lengthHex := make([]byte, 0, 8)

	for length > 0 {
		// Extract the least significant byte
		lengthHex = append(lengthHex, byte(length&0xFF))

		// Shift the number 8 bits to the right (next byte)
		length >>= 8
	}

	// Reverse the byte array
	for i, j := 0, len(lengthHex)-1; i < j; i, j = i+1, j-1 {
		lengthHex[i], lengthHex[j] = lengthHex[j], lengthHex[i]
	}

	return lengthHex
}

// EncodeBool encodes a boolean value to RLP
func EncodeBool(input bool) []byte {
	if input {
		// "true" values are encoded as 0x01
		return []byte{0x01}
	}

	// "false" values are encoded as 0x80
	return []byte{0x80}
}
