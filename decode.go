package ethrlp

import (
	"errors"
	"fmt"

	iq "github.com/madz-lab/insertion-queue"
)

var ErrInvalidLength = errors.New("invalid data length")

// DecodeBytes attempts to decode the given bytes from RLP
func DecodeBytes(input []byte) (Value, error) {
	metadata, err := getMetadata(input)
	if err != nil {
		return nil, err
	}

	var data []byte

	isListType := metadata.dataType == shortArrayType || metadata.dataType == longArrayType
	isSingleByte := metadata.dataType == byteType

	if isSingleByte {
		data = input[0:metadata.dataLength]
	} else {
		data = input[metadata.dataOffset+1 : metadata.dataLength+1]
	}

	if !isListType {
		return &BytesValue{
			value: data,
		}, nil
	}

	// Spin up the decoder worker pool
	decodePool := newDecodeWorkerPool()
	defer decodePool.close()

	decodedElements := iq.NewQueue()
	jobIndex := 0
	parseIndex := 0
	numElems := 0
	arrayLength := metadata.dataLength - metadata.dataOffset

	// Spawn a worker for each element of the array
	for parseIndex < arrayLength {
		metadata, err = getMetadata(data[parseIndex:])
		if err != nil {
			return nil, err
		}

		decodePool.addJob(&workerJob{
			storeIndex: jobIndex,
			sourceData: data[parseIndex : parseIndex+metadata.dataLength+1],
		})

		parseIndex += metadata.dataLength + 1
		jobIndex++
		numElems++
	}

	for i := 0; i < numElems; i++ {
		decodedElements.Push(decodePool.getResult())
	}

	decodedElementsLocal := make([]Value, 0, decodedElements.Len())

	for decodedElements.Len() > 0 {
		result, _ := decodedElements.PopFront().(*workerDecodeResult)
		if result.error != nil {
			return nil, fmt.Errorf("unable to decode element, %w", result.error)
		}

		decodedElementsLocal = append(decodedElementsLocal, result.decodedValue)
	}

	return &ListValue{
		values: decodedElementsLocal,
	}, nil
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
func getMetadata(data []byte) (*metadata, error) {
	if len(data) == 0 {
		return &metadata{
			dataType: emptyType,
		}, nil
	}

	firstByte := data[0]

	switch {
	case firstByte <= 0x7f:
		// Single byte value
		return &metadata{
			dataType:   byteType,
			dataOffset: 0,
			dataLength: 1,
		}, nil
	case firstByte > 0x7f && firstByte <= 0xb7:
		// Short bytes
		length := int(firstByte - 0x80)

		if length > len(data)-1 {
			return nil, constructLengthError(length, len(data)-1)
		}

		return &metadata{
			dataType:   shortBytesType,
			dataOffset: 0,
			dataLength: length,
		}, nil
	case firstByte > 0xb7 && firstByte <= 0xbf:
		// Long bytes
		lengthBytes := int(firstByte - 0xb7)
		if lengthBytes > len(data)-1 {
			return nil, constructLengthError(lengthBytes, len(data)-1)
		}

		length := convertHexArrayToInt(data[1 : lengthBytes+1])

		if length > len(data)-1-lengthBytes {
			return nil, constructLengthError(length, len(data)-1-lengthBytes)
		}

		return &metadata{
			dataType:   longBytesType,
			dataOffset: lengthBytes,
			dataLength: lengthBytes + length,
		}, nil
	case firstByte > 0xbf && firstByte <= 0xf7:
		// Short array
		length := int(firstByte - 0xc0)
		if length > len(data)-1 {
			return nil, constructLengthError(length, len(data)-1)
		}

		return &metadata{
			dataType:   shortArrayType,
			dataOffset: 0,
			dataLength: length,
		}, nil
	default:
		// Long array
		lengthBytes := int(firstByte - 0xf7)
		if lengthBytes > len(data)-1 {
			return nil, constructLengthError(lengthBytes, len(data)-1)
		}

		length := convertHexArrayToInt(data[1 : lengthBytes+1])

		if length > len(data)-1-lengthBytes {
			return nil, constructLengthError(length, len(data)-1-lengthBytes)
		}

		return &metadata{
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
