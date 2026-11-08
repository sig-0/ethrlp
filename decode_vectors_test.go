package ethrlp

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// encodeValue re-encodes a decoded tree with the package encoders.
// For canonical input it gives back the bytes that were decoded.
// It also fails the test if any node's GetValue() is a nil slice: the
// encoders map nil and empty to the same bytes, so the round trip
// alone cannot see that regression
func encodeValue(t *testing.T, value Value) []byte {
	t.Helper()

	require.NotNil(t, value.GetValue())

	if b, ok := value.GetValue().([]byte); ok {
		return EncodeBytes(b)
	}

	children := value.GetValue().([]Value)
	encoded := make([][]byte, 0, len(children))

	for _, child := range children {
		encoded = append(encoded, encodeValue(t, child))
	}

	return EncodeArray(encoded)
}

// decodeWithin runs DecodeBytes and fails the test if it does not return
// within a second, so a regression to the old infinite loop fails fast
// instead of stalling the suite until -timeout
func decodeWithin(t *testing.T, input []byte) (Value, error) {
	t.Helper()

	type result struct {
		value Value
		err   error
	}

	done := make(chan result, 1)

	go func() {
		value, err := DecodeBytes(input)
		done <- result{value: value, err: err}
	}()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	select {
	case r := <-done:
		return r.value, r.err
	case <-timer.C:
		t.Fatalf("DecodeBytes(%x) did not return within 1s", input)

		return nil, nil
	}
}

// nestedLists returns depth empty lists nested inside each other
func nestedLists(depth int) []byte {
	encoded := EncodeArray(nil)
	for i := 1; i < depth; i++ {
		encoded = EncodeArray([][]byte{encoded})
	}

	return encoded
}

func TestDecode_Valid(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name string // test name in the source
		hex  string // input, without 0x
	}{
		// ethereum/tests v17.2 (c67e485ff8b5be9abc8ad15345ec21aa22e290d9) RLPTests/rlptest.json
		{name: "emptystring", hex: "80"},
		{name: "bytestring00", hex: "00"},
		{name: "bytestring01", hex: "01"},
		{name: "bytestring7F", hex: "7f"},
		{name: "shortstring", hex: "83646f67"},
		{name: "shortstring2", hex: "b74c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c69"},
		{name: "longstring", hex: "b8384c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c6974"},
		{name: "longstring2", hex: "b904004c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e73656374657475722061646970697363696e6720656c69742e20437572616269747572206d6175726973206d61676e612c20737573636970697420736564207665686963756c61206e6f6e2c20696163756c697320666175636962757320746f72746f722e2050726f696e20737573636970697420756c74726963696573206d616c6573756164612e204475697320746f72746f7220656c69742c2064696374756d2071756973207472697374697175652065752c20756c7472696365732061742072697375732e204d6f72626920612065737420696d70657264696574206d6920756c6c616d636f7270657220616c6971756574207375736369706974206e6563206c6f72656d2e2041656e65616e2071756973206c656f206d6f6c6c69732c2076756c70757461746520656c6974207661726975732c20636f6e73657175617420656e696d2e204e756c6c6120756c74726963657320747572706973206a7573746f2c20657420706f73756572652075726e6120636f6e7365637465747572206e65632e2050726f696e206e6f6e20636f6e76616c6c6973206d657475732e20446f6e65632074656d706f7220697073756d20696e206d617572697320636f6e67756520736f6c6c696369747564696e2e20566573746962756c756d20616e746520697073756d207072696d697320696e206661756369627573206f726369206c756374757320657420756c74726963657320706f737565726520637562696c69612043757261653b2053757370656e646973736520636f6e76616c6c69732073656d2076656c206d617373612066617563696275732c2065676574206c6163696e6961206c616375732074656d706f722e204e756c6c61207175697320756c747269636965732070757275732e2050726f696e20617563746f722072686f6e637573206e69626820636f6e64696d656e74756d206d6f6c6c69732e20416c697175616d20636f6e73657175617420656e696d206174206d65747573206c75637475732c206120656c656966656e6420707572757320656765737461732e20437572616269747572206174206e696268206d657475732e204e616d20626962656e64756d2c206e6571756520617420617563746f72207472697374697175652c206c6f72656d206c696265726f20616c697175657420617263752c206e6f6e20696e74657264756d2074656c6c7573206c65637475732073697420616d65742065726f732e20437261732072686f6e6375732c206d65747573206163206f726e617265206375727375732c20646f6c6f72206a7573746f20756c747269636573206d657475732c20617420756c6c616d636f7270657220766f6c7574706174"},
		{name: "zero", hex: "80"},
		{name: "smallint", hex: "01"},
		{name: "smallint2", hex: "10"},
		{name: "smallint3", hex: "4f"},
		{name: "smallint4", hex: "7f"},
		{name: "mediumint1", hex: "8180"},
		{name: "mediumint2", hex: "8203e8"},
		{name: "mediumint3", hex: "830186a0"},
		{name: "mediumint4", hex: "8f102030405060708090a0b0c0d0e0f2"},
		{name: "mediumint5", hex: "9c0100020003000400050006000700080009000a000b000c000d000e01"},
		{name: "emptylist", hex: "c0"},
		{name: "stringlist", hex: "cc83646f6783676f6483636174"},
		{name: "multilist", hex: "c6827a77c10401"},
		{name: "shortListMax1", hex: "f784617364668471776572847a78637684617364668471776572847a78637684617364668471776572847a78637684617364668471776572"},
		{name: "longList1", hex: "f840cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376"},
		{name: "longList2", hex: "f90200" + strings.Repeat("cf84617364668471776572847a786376", 32)},
		{name: "listsoflists", hex: "c4c2c0c0c0"},
		{name: "listsoflists2", hex: "c7c0c1c0c3c0c1c0"},
		{name: "dictTest1", hex: "ecca846b6579318476616c31ca846b6579328476616c32ca846b6579338476616c33ca846b6579348476616c34"},
		{name: "bigint", hex: "a1010000000000000000000000000000000000000000000000000000000000000000"},

		// ethereum/ethereum-rlp 8c0a5af tests/test_rlp.py
		{name: "test_decode_to_bytes__single_byte_less_than_128", hex: "78"},
		{name: "test_decode_to_bytes__single_byte_gte_128", hex: "8183"},
		{name: "test_decode_to_bytes__single_byte_gte_128", hex: "8180"},
		{name: "test_decode_to_bytes__55", hex: "b7" + strings.Repeat("83", 55)},
		{name: "test_decode__to_255_uint", hex: "81ff"},
		{name: "test_decode__one_char_str", hex: "68"},
		{name: "test_decode__multi_char_str", hex: "8568656c6c6f"},
		{name: "test_decode_to_sequence__1_elem_sequence_of_byte", hex: "c68568656c6c6f"},
		{name: "test_decode_to_sequence__1_elem_sequence_of_uint", hex: "c281ff"},
		{name: "test_decode_to_sequence__10_elem_sequence_of_bytes_and_uints", hex: "e38568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f2323232323"},
		{name: "test_decode_to_sequence__20_elem_sequence_of_bytes_and_uints", hex: "f8468568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f8568656c6c6f23232323232323232323"},
		{name: "test_decode_to_sequence__nested", hex: "df8568656c6c6f81ffd683686f77d18361726583796f75c885646f696e67c123"},
		{name: "test_decode__successfully", hex: "db8568656c6c6fd483686f77cf8361726583796f75c685646f696e67"},
		{name: "test_encode_sequence__nested", hex: "dd8568656c6c6f81ffd483686f77cf8361726583796f75c685646f696e67"},
		{name: "test_encode__dataclass", hex: "c80181ffc4c38080c0"},

		// ethereum/execution-specs 16898000 tests/json_loader/test_rlp.py (valid byte strings; the integer-layer rejection of 00, 820001, 83000001 belongs to callers)
		{name: "uint 0", hex: "80"},
		{name: "uint 1", hex: "01"},
		{name: "uint 0x7f", hex: "7f"},
		{name: "uint 0x80", hex: "8180"},
		{name: "uint 0x0100", hex: "820100"},
		{name: "non-canonical uint 00", hex: "00"},
		{name: "non-canonical uint 820001", hex: "820001"},
		{name: "non-canonical uint 83000001", hex: "83000001"},

		// Yellow Paper Appendix B / ethereum.org RLP examples
		{name: "dog", hex: "83646f67"},
		{name: "cat dog", hex: "c88363617483646f67"},
		{name: "empty string", hex: "80"},
		{name: "empty list", hex: "c0"},
		{name: "integer 0", hex: "80"},
		{name: "byte 00", hex: "00"},
		{name: "byte 0f", hex: "0f"},
		{name: "integer 1024", hex: "820400"},
		{name: "set-theoretic three", hex: "c7c0c1c0c3c0c1c0"},
		{name: "lorem 56 bytes", hex: "b8384c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c6974"},

		// ethrlp: boundaries of the canonical RLP rules (55/56 lengths, nested items)
		{name: "shortest long string", hex: "b838" + strings.Repeat("00", 56)},
		{name: "shortest long list", hex: "f838" + strings.Repeat("80", 56)},
		{name: "two-byte string length", hex: "b90100" + strings.Repeat("00", 256)},
		{name: "two-byte list length", hex: "f90100" + strings.Repeat("c0", 256)},
		{name: "list of 0x80", hex: "c180"},
		{name: "list of empty list", hex: "c1c0"},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			input := hexToBytes(t, testCase.hex)

			value, err := DecodeBytes(input)
			require.NoError(t, err)

			require.Equal(t, input, encodeValue(t, value))

			// Every proper prefix is truncated, and one more byte is trailing
			for i := range input {
				_, err := DecodeBytes(input[:i])
				require.ErrorIs(t, err, ErrInvalidLength, "prefix %x", input[:i])
			}

			_, err = DecodeBytes(append(input, 0x00))
			require.ErrorIs(t, err, ErrTrailingBytes)
		})
	}
}

func TestDecode_Invalid(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		err  error  // the sentinel DecodeBytes must return
		name string // test name in the source
		hex  string // input, without 0x
	}{
		// ethereum/tests v17.2 (c67e485ff8b5be9abc8ad15345ec21aa22e290d9) RLPTests/invalidRLPTest.json
		{name: "int32Overflow", hex: "bf0f000000000000021111", err: ErrInvalidLength},
		{name: "int32Overflow2", hex: "ff0f000000000000021111", err: ErrInvalidLength},
		{name: "wrongSizeList", hex: "f80180", err: ErrNonCanonical},
		{name: "wrongSizeList2", hex: "f80100", err: ErrNonCanonical},
		{name: "incorrectLengthInArray", hex: "b9002100dc2b275d0f74e8a53e6f4ec61b27f24278820be3f82ea2110e582081b0565df0", err: ErrNonCanonical},
		{name: "randomRLP", hex: "f861f83eb9002100dc2b275d0f74e8a53e6f4ec61b27f24278820be3f82ea2110e582081b0565df027b90015002d5ef8325ae4d034df55d4b58d0dfba64d61ddd17be00000b9001a00dae30907045a2f66fa36f2bb8aa9029cbb0b8a7b3b5c435ab331", err: ErrNonCanonical},
		{name: "bytesShouldBeSingleByte00", hex: "8100", err: ErrNonCanonical},
		{name: "bytesShouldBeSingleByte01", hex: "8101", err: ErrNonCanonical},
		{name: "bytesShouldBeSingleByte7F", hex: "817F", err: ErrNonCanonical},
		{name: "emptyEncoding", hex: "", err: ErrInvalidLength},
		{name: "lessThanShortLengthArray1", hex: "81", err: ErrInvalidLength},
		{name: "lessThanShortLengthArray2", hex: "a0000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e", err: ErrInvalidLength},
		{name: "lessThanShortLengthList1", hex: "c5010203", err: ErrInvalidLength},
		{name: "lessThanShortLengthList2", hex: "e201020304050607", err: ErrInvalidLength},
		{name: "lessThanLongLengthArray1", hex: "ba010000aabbccddeeff", err: ErrInvalidLength},
		{name: "lessThanLongLengthArray2", hex: "b840ffeeddccbbaa99887766554433221100", err: ErrInvalidLength},
		{name: "lessThanLongLengthList1", hex: "f90180", err: ErrInvalidLength},
		{name: "lessThanLongLengthList2", hex: "ffffffffffffffffff0001020304050607", err: ErrInvalidLength},
		{name: "leadingZerosInLongLengthArray1", hex: "b90040000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f", err: ErrNonCanonical},
		{name: "leadingZerosInLongLengthArray2", hex: "b800", err: ErrNonCanonical},
		{name: "leadingZerosInLongLengthList1", hex: "fb00000040000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f", err: ErrNonCanonical},
		{name: "leadingZerosInLongLengthList2", hex: "f800", err: ErrNonCanonical},
		{name: "nonOptimalLongLengthArray1", hex: "b81000112233445566778899aabbccddeeff", err: ErrNonCanonical},
		{name: "nonOptimalLongLengthArray2", hex: "b801ff", err: ErrNonCanonical},
		{name: "nonOptimalLongLengthList1", hex: "f810000102030405060708090a0b0c0d0e0f", err: ErrNonCanonical},
		{name: "nonOptimalLongLengthList2", hex: "f803112233", err: ErrNonCanonical},

		// ethrlp v0.2.1 regressions: inputs the old decoder hung on, panicked on or accepted
		{name: "hang: negative string length in list", hex: "c9bffffffffffffffff7", err: ErrInvalidLength},
		{name: "hang: negative list length in list", hex: "c9fffffffffffffffff7", err: ErrInvalidLength},
		{name: "panic: negative slice bound", hex: "ff8f3030303030303030", err: ErrInvalidLength},
		{name: "trailing byte after list", hex: "c20102ff", err: ErrTrailingBytes},

		// ethereum/ethereum-rlp 8c0a5af tests/test_rlp.py
		{name: "test_decode_to_bytes__out_of_bounds", hex: "84", err: ErrInvalidLength},
		{name: "test_decode_to_bytes__single_byte_sequence", hex: "8179", err: ErrNonCanonical},
		{name: "test_decode_to_bytes__long_missing_length", hex: "b8", err: ErrInvalidLength},
		{name: "test_decode_to_bytes__long_out_of_bounds", hex: "ba370000", err: ErrInvalidLength},
		{name: "test_decode_to_bytes__rejects_trailing_bytes", hex: "8000", err: ErrTrailingBytes},
		{name: "test_decode_to_bytes__rejects_trailing_bytes_long", hex: "b838" + strings.Repeat("00", 56) + "00", err: ErrTrailingBytes},
		{name: "test_decode_to_bytes__negative_length", hex: "7ebcc55ebeff", err: ErrTrailingBytes},
		{name: "test_decode_to_bytes__negative_length", hex: "5135393abaf4da05b7", err: ErrTrailingBytes},
		{name: "test_decode_to_sequence__out_of_bounds_item", hex: "c181", err: ErrInvalidLength},
		{name: "test_decode_to_sequence__rejects_trailing_bytes", hex: "c000", err: ErrTrailingBytes},
		{name: "test_decode_to_sequence__rejects_trailing_bytes_long", hex: "f838" + strings.Repeat("80", 56) + "00", err: ErrTrailingBytes},
		{name: "test_decode_to_sequence__out_of_bounds", hex: "c1", err: ErrInvalidLength},
		{name: "test_decode_to_sequence__long_missing_length", hex: "f8", err: ErrInvalidLength},
		{name: "test_decode_to_sequence__long_zero", hex: "f800", err: ErrNonCanonical},
		{name: "test_decode_to_sequence__long_too_short", hex: "f837" + strings.Repeat("00", 55), err: ErrNonCanonical},
		{name: "test_decode_to_sequence__long_out_of_bounds", hex: "f839", err: ErrInvalidLength},
		{name: "test_decode_item_length__long_length_out_of_bounds", hex: "bf01000000000000", err: ErrInvalidLength},
		{name: "test_decode_item_length__long", hex: "bf0100000000000000", err: ErrInvalidLength},
		{name: "test_decode_item_length__long_length_too_short", hex: "bf0000000000000000", err: ErrNonCanonical},
		{name: "test_decode_item_length__long_long_length_too_short", hex: "f80000000000000000", err: ErrNonCanonical},
		{name: "test_decode_item_length__long_long_length_out_of_bounds", hex: "f901", err: ErrInvalidLength},
		{name: "test_decode_item_length__long_long", hex: "f8ff", err: ErrInvalidLength},

		// ethrlp: boundaries of the canonical RLP rules (55/56 lengths, nested items)
		{name: "longest short form in long form", hex: "b837" + strings.Repeat("00", 55), err: ErrNonCanonical},
		{name: "non-canonical byte inside a list", hex: "c28101", err: ErrNonCanonical},
		{name: "long form under 56 inside a list", hex: "c3b80100", err: ErrNonCanonical},
		{name: "child overruns its list", hex: "c28361626364", err: ErrInvalidLength},
		{name: "child header overruns its list", hex: "c1c1", err: ErrInvalidLength},
		{name: "list overruns input", hex: "c3c2c0", err: ErrInvalidLength},
		{name: "truncated two-byte length", hex: "b901", err: ErrInvalidLength},
		{name: "largest 8-byte length", hex: "bfffffffffffffffff", err: ErrInvalidLength},
		{name: "long string without payload", hex: "b838", err: ErrInvalidLength},
		{name: "long list without payload", hex: "f838", err: ErrInvalidLength},
		{name: "two-byte list length without payload", hex: "f90100", err: ErrInvalidLength},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			value, err := decodeWithin(t, hexToBytes(t, testCase.hex))
			require.ErrorIs(t, err, testCase.err)
			require.Nil(t, value)
		})
	}

	t.Run("nil input", func(t *testing.T) {
		t.Parallel()

		value, err := decodeWithin(t, nil)
		require.ErrorIs(t, err, ErrInvalidLength)
		require.Nil(t, value)
	})
}

func TestDecode_AllShortInputs(t *testing.T) {
	t.Parallel()

	accepted := func(input []byte) bool {
		value, err := DecodeBytes(input)
		if err != nil {
			return false
		}

		require.Equal(t, input, encodeValue(t, value))

		return true
	}

	oneByte := 0

	for b := 0; b <= 0xff; b++ {
		if accepted([]byte{byte(b)}) {
			oneByte++
		}
	}

	twoBytes := 0

	for b := 0; b <= 0xffff; b++ {
		if accepted([]byte{byte(b >> 8), byte(b)}) {
			twoBytes++
		}
	}

	// one byte: 00..7f, 80, c0 (130); two bytes: 81 80..ff, c1 00..7f, c180, c1c0 (258)
	require.Equal(t, 130, oneByte)
	require.Equal(t, 258, twoBytes)
}

func TestDecode_DeepNesting(t *testing.T) {
	t.Parallel()

	// RLP has no nesting limit, so deeply nested lists decode
	input := nestedLists(10_000)

	value, err := DecodeBytes(input)
	require.NoError(t, err)
	require.Equal(t, input, encodeValue(t, value))
}

func TestDecode_LargeInputs(t *testing.T) {
	t.Parallel()

	// ethereum/ethereum-rlp 8c0a5af tests/test_rlp.py
	payload := bytes.Repeat([]byte{0x83}, 1<<20)

	t.Run("test_decode_to_bytes__large", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(append([]byte{0xba, 0x10, 0x00, 0x00}, payload...))
		require.NoError(t, err)
		require.Equal(t, payload, value.GetValue())
	})

	t.Run("test_decode_to_bytes__long_zero", func(t *testing.T) {
		t.Parallel()

		_, err := DecodeBytes(append([]byte{0xba, 0x00, 0x00, 0x00}, payload...))
		require.ErrorIs(t, err, ErrNonCanonical)
	})

	t.Run("test_decode_to_bytes__long_too_short", func(t *testing.T) {
		t.Parallel()

		_, err := DecodeBytes(append([]byte{0xb8, 0x37, 0x00, 0x00}, payload...))
		require.ErrorIs(t, err, ErrNonCanonical)
	})

	t.Run("test_decode_joined_encodings__round_trips_multi_byte_length_prefixes", func(t *testing.T) {
		t.Parallel()

		for _, n := range []int{55, 56, 255, 256, 65535, 65536} {
			input := EncodeArray([][]byte{
				EncodeBytes([]byte{0x11}),
				EncodeBytes(bytes.Repeat([]byte{0x22}, n)),
				EncodeArray([][]byte{
					EncodeBytes(bytes.Repeat([]byte{0x33}, n)),
					EncodeBytes([]byte{0x44}),
				}),
				EncodeArray([][]byte{
					EncodeArray([][]byte{EncodeBytes(bytes.Repeat([]byte{0x55}, n))}),
				}),
				EncodeBytes(bytes.Repeat([]byte{0x66}, n)),
			})

			value, err := DecodeBytes(input)
			require.NoError(t, err, "size %d", n)
			require.Equal(t, input, encodeValue(t, value), "size %d", n)
		}
	})
}
