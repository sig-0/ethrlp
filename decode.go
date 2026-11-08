package ethrlp

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidLength = errors.New("invalid data length")
	ErrNonCanonical  = errors.New("non-canonical RLP")
	ErrTrailingBytes = errors.New("trailing bytes after RLP item")
)

// DecodeBytes attempts to decode the given bytes from RLP
func DecodeBytes(input []byte) (Value, error) {
	// Decode the top-level item
	value, rest, err := decodeItem(input)
	if err != nil {
		return nil, err
	}

	// Make sure there is no data after the top-level item
	if len(rest) != 0 {
		return nil, fmt.Errorf("%w: %dB", ErrTrailingBytes, len(rest))
	}

	return value, nil
}

// decodeItem decodes a single RLP item, and returns the bytes that follow it
func decodeItem(input []byte) (Value, []byte, error) {
	if len(input) == 0 {
		return nil, nil, constructLengthError(1, 0)
	}

	// Single byte value
	if input[0] < 0x80 {
		return BytesValue{value: input[:1]}, input[1:], nil
	}

	// Fetch the header size and the data length
	headerSize, length, err := parseHeader(input)
	if err != nil {
		return nil, nil, err
	}

	rest := input[headerSize:]

	// Check range bounds
	if length > uint64(len(rest)) {
		return nil, nil, constructLengthError(length, len(rest))
	}

	// Extract the payload bytes that belong to this RLP item
	payload, rest := rest[:length], rest[length:]

	// If it’s not a list, simply return BytesValue (byte slice)
	if input[0] < 0xc0 {
		// A single byte below 0x80 must be encoded as itself
		if len(payload) == 1 && payload[0] < 0x80 {
			return nil, nil, fmt.Errorf("%w: byte %#x must be encoded as itself", ErrNonCanonical, payload[0])
		}

		return BytesValue{value: payload}, rest, nil
	}

	values := make([]Value, 0, 9) // 9 is chosen as a common Ethereum RLP list size

	// Parse each element of the list
	for len(payload) > 0 {
		var value Value

		// Decode the item recursively, and move to the next element
		value, payload, err = decodeItem(payload)
		if err != nil {
			return nil, nil, err
		}

		// Save the decoded item
		values = append(values, value)
	}

	return ListValue{values: values}, rest, nil
}

// parseHeader returns the header size and the data length of a bytes or list item
func parseHeader(input []byte) (int, uint64, error) {
	// Bytes header bases
	shortBase, longBase := byte(0x80), byte(0xb7)

	if input[0] >= 0xc0 {
		// List header bases
		shortBase, longBase = 0xc0, 0xf7
	}

	// Short bytes / short array
	if input[0] <= longBase {
		return 1, uint64(input[0] - shortBase), nil
	}

	// Long bytes / long array
	lengthSize := int(input[0] - longBase)
	if lengthSize > len(input)-1 {
		return 0, 0, constructLengthError(uint64(input[0]-longBase), len(input)-1)
	}

	// The length cannot have leading zeros
	if input[1] == 0 {
		return 0, 0, fmt.Errorf("%w: length has a leading zero byte", ErrNonCanonical)
	}

	var length uint64

	for _, b := range input[1 : 1+lengthSize] {
		// Shift the current length value 8 bits to the left,
		// and add the current byte to the length
		length = length<<8 | uint64(b)
	}

	// Lengths under 56 must use the short form
	if length < 56 {
		return 0, 0, fmt.Errorf("%w: long form for length %d, under 56", ErrNonCanonical, length)
	}

	return 1 + lengthSize, length, nil
}

// constructLengthError constructs an invalid RLP length error
func constructLengthError(expected uint64, actual int) error {
	return fmt.Errorf(
		"%w: expected %dB, got %dB",
		ErrInvalidLength,
		expected,
		actual,
	)
}
