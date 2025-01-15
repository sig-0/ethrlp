package ethrlp

import (
	"errors"
	"fmt"
)

var ErrInvalidLength = errors.New("invalid data length")

// DecodeBytes attempts to decode the given bytes from RLP
func DecodeBytes(input []byte) (Value, error) {
	// Fetch the top-level metadata
	topMeta, err := getMetadata(input)
	if err != nil {
		return nil, err
	}

	// Detect whether this is a list or a single byte or something else
	var (
		isListType   = topMeta.dataType == shortArrayType || topMeta.dataType == longArrayType
		isSingleByte = topMeta.dataType == byteType
	)

	// Extract the payload bytes that belong to this RLP item
	var data []byte
	if isSingleByte {
		// For a single-byte item, data length should be 1
		data = input[:topMeta.dataLength]
	} else {
		// Otherwise, skip the first byte (+ any length-bytes),
		// then take <data length> bytes
		data = input[topMeta.dataOffset+1 : topMeta.dataLength+1]
	}

	// If it’s not a list, simply return BytesValue (byte slice)
	if !isListType {
		return BytesValue{value: data}, nil
	}

	var (
		// Calculate how many data bytes we have in this list
		// (excludes the prefix & length bytes)
		listLength   = topMeta.dataLength - topMeta.dataOffset
		decodedItems = make([]Value, 0, 9) // 9 is chosen as a common Ethereum RLP list size
	)

	// Parse each element of the list
	for parseIndex := 0; parseIndex < listLength; {
		// Get metadata on the element
		elemMeta, err := getMetadata(data[parseIndex:])
		if err != nil {
			return nil, err
		}

		// Calculate the total on-wire size of this element.
		// A single byte in [0x00..0x7f] is its *own* entire encoding (no extra prefix)
		itemTotal := 1
		if elemMeta.dataType != byteType {
			// For bytes / lists, the total is prefix byte(s) + data length
			itemTotal = 1 + elemMeta.dataLength
		}

		// Check range bounds
		if parseIndex+itemTotal > len(data) {
			return nil, fmt.Errorf(
				"RLP data is truncated: parse index =%d total items=%d data length=%d",
				parseIndex, itemTotal, len(data),
			)
		}

		// Extract the sub-slice for this item
		itemBytes := data[parseIndex : parseIndex+itemTotal]

		// Decode the item recursively
		decodedItem, err := DecodeBytes(itemBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to decode item, %w", err)
		}

		// Save the decoded item
		decodedItems = append(decodedItems, decodedItem)

		// Move to the next element
		parseIndex += itemTotal
	}

	return ListValue{values: decodedItems}, nil
}

const (
	emptyType = iota
	byteType
	shortBytesType
	longBytesType
	shortArrayType
	longArrayType
)

type metadata struct {
	dataType   int // type of data
	dataOffset int // where the data starts (not including first byte)
	dataLength int // total data size (not including first byte)
}

// getMetadata returns the metadata about the top-level RLP type
func getMetadata(data []byte) (metadata, error) {
	if len(data) == 0 {
		return metadata{
			dataType: emptyType,
		}, nil
	}

	firstByte := data[0]

	switch {
	case firstByte <= 0x7f:
		// Single byte value
		return metadata{
			dataType:   byteType,
			dataOffset: 0,
			dataLength: 1,
		}, nil
	case firstByte > 0x7f && firstByte <= 0xb7:
		// Short bytes
		length := int(firstByte - 0x80)

		if length > len(data)-1 {
			return metadata{}, constructLengthError(length, len(data)-1)
		}

		return metadata{
			dataType:   shortBytesType,
			dataOffset: 0,
			dataLength: length,
		}, nil
	case firstByte > 0xb7 && firstByte <= 0xbf:
		// Long bytes
		lengthBytes := int(firstByte - 0xb7)
		if lengthBytes > len(data)-1 {
			return metadata{}, constructLengthError(lengthBytes, len(data)-1)
		}

		length := convertHexArrayToInt(data[1 : lengthBytes+1])

		if length > len(data)-1-lengthBytes {
			return metadata{}, constructLengthError(length, len(data)-1-lengthBytes)
		}

		return metadata{
			dataType:   longBytesType,
			dataOffset: lengthBytes,
			dataLength: lengthBytes + length,
		}, nil
	case firstByte > 0xbf && firstByte <= 0xf7:
		// Short array
		length := int(firstByte - 0xc0)
		if length > len(data)-1 {
			return metadata{}, constructLengthError(length, len(data)-1)
		}

		return metadata{
			dataType:   shortArrayType,
			dataOffset: 0,
			dataLength: length,
		}, nil
	default:
		// Long array
		lengthBytes := int(firstByte - 0xf7)
		if lengthBytes > len(data)-1 {
			return metadata{}, constructLengthError(lengthBytes, len(data)-1)
		}

		length := convertHexArrayToInt(data[1 : lengthBytes+1])

		if length > len(data)-1-lengthBytes {
			return metadata{}, constructLengthError(length, len(data)-1-lengthBytes)
		}

		return metadata{
			dataType:   longArrayType,
			dataOffset: lengthBytes,
			dataLength: lengthBytes + length,
		}, nil
	}
}

// convertHexArrayToInt converts the byte array of hex values
// to its corresponding integer representation
func convertHexArrayToInt(hexArray []byte) int {
	length := 0

	for _, b := range hexArray {
		// Shift the current length value 8 bits to the left
		length <<= 8

		// Add the current byte to the length
		length |= int(b)
	}

	return length
}

// constructLengthError constructs an invalid RLP length error
func constructLengthError(expected, actual int) error {
	return fmt.Errorf(
		"%w: expected %dB, got %dB",
		ErrInvalidLength,
		expected,
		actual,
	)
}
