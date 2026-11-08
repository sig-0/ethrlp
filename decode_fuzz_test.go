package ethrlp

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzDecodeBytes(f *testing.F) {
	// Representative seeds: valid strings and lists of each length form, and one
	// input per rejection rule. The full vector tables are in TestDecode_Valid
	// and TestDecode_Invalid
	seeds := []string{
		"80", "00", "7f", "8180", "8203e8", // strings: empty, single bytes, short (integer 1000)
		"b838" + strings.Repeat("00", 56),        // shortest long-form string
		"c0", "c88363617483646f67", "c4c2c0c0c0", // lists: empty, ["cat", "dog"], nested empty lists
		"f838" + strings.Repeat("80", 56),  // shortest long-form list
		"8101", "b800", "f80180", "c28101", // non-canonical
		"8000", "c20102ff", // trailing bytes
		"81", "c5010203", "bf0f000000000000021111", // truncated or overflowing lengths
		"c9bffffffffffffffff7", "c9fffffffffffffffff7", "ff8f3030303030303030", // v0.2.1 hang and panic inputs
	}

	for _, seed := range seeds {
		f.Add(hexToBytes(f, seed))
	}

	f.Add([]byte(nil))

	f.Fuzz(func(t *testing.T, input []byte) {
		value, err := decodeWithin(t, input)
		if err != nil {
			require.Nil(t, value)
			require.True(t,
				errors.Is(err, ErrInvalidLength) ||
					errors.Is(err, ErrNonCanonical) ||
					errors.Is(err, ErrTrailingBytes),
				"unexpected error %v for %x", err, input,
			)

			return
		}

		require.Equal(t, input, encodeValue(t, value))
	})
}
