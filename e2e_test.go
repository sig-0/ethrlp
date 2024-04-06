package ethrlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode_Byte(t *testing.T) {
	t.Parallel()

	t.Run("Byte in [0x00, 0x7f] range", func(t *testing.T) {
		t.Parallel()

		for i := 0x00; i <= 0x7f; i++ {
			if i == 0x1 {
				// Skip true boolean value
				continue
			}

			// Encode and decode the value
			decoded, err := DecodeBytes(EncodeByte(byte(i)))

			require.NoError(t, err)
			require.NotNil(t, decoded)

			// Make sure the type is valid
			assert.Equal(t, Bytes, decoded.GetType())

			// Make sure the underlying Go type
			// is valid
			castValue, ok := decoded.GetValue().([]byte)
			assert.True(t, ok)

			// Make sure the decoded value
			// matches the raw value
			assert.Equal(t, byte(i), castValue[0])
		}
	})

	t.Run("Byte in the [0x81, 0xff] range", func(t *testing.T) {
		t.Parallel()

		for i := 0x81; i <= 0xff; i++ {
			// Encode and decode the value
			decoded, err := DecodeBytes(EncodeByte(byte(i)))

			require.NoError(t, err)
			require.NotNil(t, decoded)

			// Make sure the type is valid
			assert.Equal(t, Bytes, decoded.GetType())

			// Make sure the underlying Go type
			// is valid
			castValue, ok := decoded.GetValue().([]byte)
			assert.True(t, ok)

			// Make sure the decoded value
			// matches the raw value
			assert.Equal(t, byte(i), castValue[0])
		}
	})
}

func TestEncodeDecode_Bytes(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name  string
		input []byte
	}{
		{
			"Input bytes <=55B in length",
			[]byte("hello world"),
		},
		{
			"Input bytes >55B in length",
			[]byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit"),
		},
		{
			"Input bytes >55B in length, long",
			[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat"),
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Encode and decode the value
			decoded, err := DecodeBytes(EncodeBytes(testCase.input))

			require.NoError(t, err)
			require.NotNil(t, decoded)

			// Make sure the type is valid
			assert.Equal(t, Bytes, decoded.GetType())

			// Make sure the underlying Go type
			// is valid
			castValue, ok := decoded.GetValue().([]byte)
			assert.True(t, ok)

			// Make sure the decoded value
			// matches the raw value
			assert.Equal(t, testCase.input, castValue)
		})
	}
}

func TestEncodeDecode_Array_Bytes(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name  string
		input [][]byte
	}{
		{
			"Empty array",
			[][]byte{},
		},
		{
			"Populated array <=55B in length",
			[][]byte{
				EncodeString("hello"),
				EncodeString("world"),
			},
		},
		{
			"Populated array >55B in length",
			[][]byte{
				EncodeString("aaa"),
				EncodeString("bbb"),
				EncodeString("ccc"),
				EncodeString("ddd"),
				EncodeString("eee"),
				EncodeString("fff"),
				EncodeString("ggg"),
				EncodeString("hhh"),
				EncodeString("iii"),
				EncodeString("jjj"),
				EncodeString("kkk"),
				EncodeString("lll"),
				EncodeString("mmm"),
				EncodeString("nnn"),
				EncodeString("ooo"),
			},
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Encode and decode the value
			decoded, err := DecodeBytes(EncodeArray(testCase.input))

			require.NoError(t, err)
			require.NotNil(t, decoded)

			// Make sure the type is valid
			assert.Equal(t, List, decoded.GetType())

			// Make sure the underlying Go type is valid
			castValue, ok := decoded.GetValue().([]Value)
			assert.True(t, ok)

			// Make sure the decoded value
			// matches the raw value
			require.Len(t, castValue, len(testCase.input))

			for index, originalValue := range testCase.input {
				rawValue := castValue[index]

				// Make sure the type is valid
				assert.Equal(t, Bytes, rawValue.GetType())

				// Make sure the underlying Go type is valid
				arrayValue, ok := rawValue.GetValue().([]byte)
				require.True(t, ok)

				// Make sure the encoded value matches the original
				assert.Equal(t, originalValue, EncodeBytes(arrayValue))
			}
		})
	}
}

func TestEncodeDecode_Array_Nested(t *testing.T) {
	t.Parallel()

	input := [][]byte{
		EncodeArray([][]byte{
			EncodeString("key1"),
			EncodeString("val1"),
		}),
		EncodeArray([][]byte{
			EncodeString("key2"),
			EncodeString("val2"),
		}),
		EncodeArray([][]byte{
			EncodeString("key3"),
			EncodeString("val3"),
		}),
		EncodeArray([][]byte{
			EncodeString("key4"),
			EncodeString("val4"),
		}),
	}

	// Encode and decode the value
	decoded, err := DecodeBytes(EncodeArray(input))

	require.NoError(t, err)
	require.NotNil(t, decoded)

	// Make sure the type is valid
	assert.Equal(t, List, decoded.GetType())

	// Make sure the underlying Go type
	// is valid
	castValue, ok := decoded.GetValue().([]Value)
	assert.True(t, ok)

	// Make sure the decoded value
	// matches the raw value
	require.Len(t, castValue, len(input))

	for index := range input {
		rawValue := castValue[index]

		// Make sure the type is valid
		assert.Equal(t, List, rawValue.GetType())

		// Make sure the underlying Go type is valid
		nestedValues, ok := rawValue.GetValue().([]Value)
		require.True(t, ok)

		for _, nestedValue := range nestedValues {
			arrayValue, ok := nestedValue.GetValue().([]byte)
			require.True(t, ok)

			assert.NotNil(t, arrayValue)
		}
	}
}
