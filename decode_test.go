package ethrlp

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode_EthereumTests(t *testing.T) {
	t.Parallel()

	t.Run("emptystring", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes([]byte{0x80})
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte(""), castValue)
	})

	t.Run("bytestring00", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes([]byte{0x00})
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte{0x00}, castValue)
	})

	t.Run("bytestring01", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes([]byte{0x01})
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte{0x01}, castValue)
	})

	t.Run("bytestring7F", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes([]byte{0x7F})
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte{0x7F}, castValue)
	})

	t.Run("shortstring", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "83646f67"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte("dog"), castValue)
	})

	t.Run("shortstring2", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "b74c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c69"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte("Lorem ipsum dolor sit amet, consectetur adipisicing eli"), castValue)
	})

	t.Run("longstring", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "b8384c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c6974"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit"), castValue)
	})

	t.Run("longstring2", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "b904004c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e73656374657475722061646970697363696e6720656c69742e20437572616269747572206d6175726973206d61676e612c20737573636970697420736564207665686963756c61206e6f6e2c20696163756c697320666175636962757320746f72746f722e2050726f696e20737573636970697420756c74726963696573206d616c6573756164612e204475697320746f72746f7220656c69742c2064696374756d2071756973207472697374697175652065752c20756c7472696365732061742072697375732e204d6f72626920612065737420696d70657264696574206d6920756c6c616d636f7270657220616c6971756574207375736369706974206e6563206c6f72656d2e2041656e65616e2071756973206c656f206d6f6c6c69732c2076756c70757461746520656c6974207661726975732c20636f6e73657175617420656e696d2e204e756c6c6120756c74726963657320747572706973206a7573746f2c20657420706f73756572652075726e6120636f6e7365637465747572206e65632e2050726f696e206e6f6e20636f6e76616c6c6973206d657475732e20446f6e65632074656d706f7220697073756d20696e206d617572697320636f6e67756520736f6c6c696369747564696e2e20566573746962756c756d20616e746520697073756d207072696d697320696e206661756369627573206f726369206c756374757320657420756c74726963657320706f737565726520637562696c69612043757261653b2053757370656e646973736520636f6e76616c6c69732073656d2076656c206d617373612066617563696275732c2065676574206c6163696e6961206c616375732074656d706f722e204e756c6c61207175697320756c747269636965732070757275732e2050726f696e20617563746f722072686f6e637573206e69626820636f6e64696d656e74756d206d6f6c6c69732e20416c697175616d20636f6e73657175617420656e696d206174206d65747573206c75637475732c206120656c656966656e6420707572757320656765737461732e20437572616269747572206174206e696268206d657475732e204e616d20626962656e64756d2c206e6571756520617420617563746f72207472697374697175652c206c6f72656d206c696265726f20616c697175657420617263752c206e6f6e20696e74657264756d2074656c6c7573206c65637475732073697420616d65742065726f732e20437261732072686f6e6375732c206d65747573206163206f726e617265206375727375732c20646f6c6f72206a7573746f20756c747269636573206d657475732c20617420756c6c616d636f7270657220766f6c7574706174"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat"), castValue)
	})

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(0)

		value, err := DecodeBytes(hexToBytes(t, "80"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("smallint", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(1)

		value, err := DecodeBytes(hexToBytes(t, "01"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("smallint2", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(16)

		value, err := DecodeBytes(hexToBytes(t, "10"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("smallint3", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(79)

		value, err := DecodeBytes(hexToBytes(t, "4f"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("smallint4", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(127)

		value, err := DecodeBytes(hexToBytes(t, "7f"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("mediumint1", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(128)

		value, err := DecodeBytes(hexToBytes(t, "8180"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("mediumint2", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(1000)

		value, err := DecodeBytes(hexToBytes(t, "8203e8"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("mediumint3", func(t *testing.T) {
		t.Parallel()

		expectedValue := big.NewInt(100000)

		value, err := DecodeBytes(hexToBytes(t, "830186a0"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("mediumint4", func(t *testing.T) {
		t.Parallel()

		expectedValue, ok := big.NewInt(0).SetString("83729609699884896815286331701780722", 10)
		require.True(t, ok)

		value, err := DecodeBytes(hexToBytes(t, "8f102030405060708090a0b0c0d0e0f2"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("mediumint5", func(t *testing.T) {
		t.Parallel()

		expectedValue, ok := big.NewInt(0).SetString("105315505618206987246253880190783558935785933862974822347068935681", 10)
		require.True(t, ok)

		value, err := DecodeBytes(hexToBytes(t, "9c0100020003000400050006000700080009000a000b000c000d000e01"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		castValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			expectedValue.String(),
			big.NewInt(0).SetBytes(castValue).String(),
		)
	})

	t.Run("emptylist", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "c0"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValue, ok := value.GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, listValue, 0)
	})

	t.Run("stringlist", func(t *testing.T) {
		t.Parallel()

		expectedValues := [][]byte{
			[]byte("dog"),
			[]byte("god"),
			[]byte("cat"),
		}

		value, err := DecodeBytes(hexToBytes(t, "cc83646f6783676f6483636174"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValue, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValue, len(expectedValues))

		for index, v := range listValue {
			assert.Equal(t, v.GetType(), Bytes)

			rawValue, ok := v.GetValue().([]byte)
			require.True(t, ok)

			assert.Equal(t, expectedValues[index], rawValue)
		}
	})

	t.Run("multilist", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "c6827a77c10401"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValues, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValues, 3)

		// Validate the first string value "zw"
		require.Equal(t, Bytes, listValues[0].GetType())

		firstValue, ok := listValues[0].GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, []byte("zw"), firstValue)

		// Validate the second array value with a single number 4
		require.Equal(t, List, listValues[1].GetType())

		secondValue, ok := listValues[1].GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, secondValue, 1)

		require.Equal(t, Bytes, secondValue[0].GetType())

		assert.Equal(t, big.NewInt(4).Bytes(), secondValue[0].GetValue())

		// Validate the third number value 1
		require.Equal(t, Bytes, listValues[2].GetType())

		thirdValue, ok := listValues[2].GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(
			t,
			big.NewInt(1).Bytes(),
			big.NewInt(0).SetBytes(thirdValue).Bytes(),
		)
	})

	t.Run("shortListMax1", func(t *testing.T) {
		t.Parallel()

		expectedValues := [][]byte{
			[]byte("asdf"),
			[]byte("qwer"),
			[]byte("zxcv"),
			[]byte("asdf"),
			[]byte("qwer"),
			[]byte("zxcv"),
			[]byte("asdf"),
			[]byte("qwer"),
			[]byte("zxcv"),
			[]byte("asdf"),
			[]byte("qwer"),
		}

		value, err := DecodeBytes(hexToBytes(t, "f784617364668471776572847a78637684617364668471776572847a78637684617364668471776572847a78637684617364668471776572"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValue, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValue, len(expectedValues))

		for index, v := range listValue {
			assert.Equal(t, v.GetType(), Bytes)

			rawValue, ok := v.GetValue().([]byte)
			require.True(t, ok)

			assert.Equal(t, expectedValues[index], rawValue)
		}
	})

	t.Run("longList1", func(t *testing.T) {
		t.Parallel()

		expectedValues := [][]byte{
			[]byte("asdf"),
			[]byte("qwer"),
			[]byte("zxcv"),
		}

		validateNesting := func(value Value) {
			require.Equal(t, List, value.GetType())

			listValues, ok := value.GetValue().([]Value)
			require.True(t, ok)

			require.Len(t, listValues, 3)

			for index, v := range listValues {
				assert.Equal(t, v.GetType(), Bytes)

				rawValue, ok := v.GetValue().([]byte)
				require.True(t, ok)

				assert.Equal(t, expectedValues[index], rawValue)
			}
		}

		value, err := DecodeBytes(hexToBytes(t, "f840cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValues, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValues, 4)

		for _, v := range listValues {
			validateNesting(v)
		}
	})

	t.Run("longList2", func(t *testing.T) {
		t.Parallel()

		expectedValues := [][]byte{
			[]byte("asdf"),
			[]byte("qwer"),
			[]byte("zxcv"),
		}

		validateNesting := func(value Value) {
			require.Equal(t, List, value.GetType())

			listValues, ok := value.GetValue().([]Value)
			require.True(t, ok)

			require.Len(t, listValues, 3)

			for index, v := range listValues {
				assert.Equal(t, v.GetType(), Bytes)

				rawValue, ok := v.GetValue().([]byte)
				require.True(t, ok)

				assert.Equal(t, expectedValues[index], rawValue)
			}
		}

		value, err := DecodeBytes(hexToBytes(t, "f90200cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValues, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValues, 32)

		for _, v := range listValues {
			validateNesting(v)
		}
	})

	t.Run("listsoflists", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "c4c2c0c0c0"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValues, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValues, 2)

		require.Equal(t, List, listValues[0].GetType())

		firstValue, ok := listValues[0].GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, firstValue, 2)

		firstSubValue, ok := firstValue[0].GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, firstSubValue, 0)

		secondSubValue, ok := firstValue[1].GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, secondSubValue, 0)

		require.Equal(t, List, listValues[1].GetType())

		secondValue, ok := listValues[1].GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, secondValue, 0)
	})

	t.Run("listsoflists2", func(t *testing.T) {
		t.Parallel()

		value, err := DecodeBytes(hexToBytes(t, "c7c0c1c0c3c0c1c0"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValues, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValues, 3)

		// Validate the first array value
		require.Equal(t, List, listValues[0].GetType())

		firstValue, ok := listValues[0].GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, firstValue, 0)

		// Validate the second array value
		require.Equal(t, List, listValues[1].GetType())

		secondValue, ok := listValues[1].GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, secondValue, 1)

		assert.Equal(t, List, secondValue[0].GetType())

		secondSubValue, ok := secondValue[0].GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, secondSubValue, 0)

		// Validate the third array value
		thirdValue, ok := listValues[2].GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, thirdValue, 2)

		assert.Equal(t, List, thirdValue[0].GetType())

		thirdSubValueA, ok := thirdValue[0].GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, thirdSubValueA, 0)

		assert.Equal(t, List, thirdValue[1].GetType())

		thirdSubValueB, ok := thirdValue[1].GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, thirdSubValueB, 1)

		assert.Equal(t, List, thirdSubValueB[0].GetType())

		thirdSubValueBA, ok := thirdSubValueB[0].GetValue().([]Value)
		require.True(t, ok)

		assert.Len(t, thirdSubValueBA, 0)
	})

	t.Run("dictTest1", func(t *testing.T) {
		t.Parallel()

		validateNesting := func(index int, value Value) {
			require.Equal(t, List, value.GetType())

			listValues, ok := value.GetValue().([]Value)
			require.True(t, ok)

			require.Len(t, listValues, 2)

			require.Equal(t, Bytes, listValues[0].GetType())

			keyValue, ok := listValues[0].GetValue().([]byte)
			require.True(t, ok)

			assert.Equal(t, fmt.Sprintf("key%d", index+1), string(keyValue))

			require.Equal(t, Bytes, listValues[0].GetType())

			valueValue, ok := listValues[1].GetValue().([]byte)
			require.True(t, ok)

			assert.Equal(t, fmt.Sprintf("val%d", index+1), string(valueValue))
		}

		value, err := DecodeBytes(hexToBytes(t, "ecca846b6579318476616c31ca846b6579328476616c32ca846b6579338476616c33ca846b6579348476616c34"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValues, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValues, 4)

		// Validate each dictionary value
		for index, value := range listValues {
			validateNesting(index, value)
		}
	})

	t.Run("dictTest1", func(t *testing.T) {
		t.Parallel()

		validateNesting := func(index int, value Value) {
			require.Equal(t, List, value.GetType())

			listValues, ok := value.GetValue().([]Value)
			require.True(t, ok)

			require.Len(t, listValues, 2)

			require.Equal(t, Bytes, listValues[0].GetType())

			keyValue, ok := listValues[0].GetValue().([]byte)
			require.True(t, ok)

			assert.Equal(t, fmt.Sprintf("key%d", index+1), string(keyValue))

			require.Equal(t, Bytes, listValues[0].GetType())

			valueValue, ok := listValues[1].GetValue().([]byte)
			require.True(t, ok)

			assert.Equal(t, fmt.Sprintf("val%d", index+1), string(valueValue))
		}

		value, err := DecodeBytes(hexToBytes(t, "ecca846b6579318476616c31ca846b6579328476616c32ca846b6579338476616c33ca846b6579348476616c34"))
		require.NoError(t, err)

		assert.Equal(t, List, value.GetType())

		listValues, ok := value.GetValue().([]Value)
		require.True(t, ok)

		require.Len(t, listValues, 4)

		// Validate each dictionary value
		for index, value := range listValues {
			validateNesting(index, value)
		}
	})

	t.Run("bigint", func(t *testing.T) {
		t.Parallel()

		expectedNumber, ok := big.NewInt(0).SetString(
			"115792089237316195423570985008687907853269984665640564039457584007913129639936",
			10,
		)
		require.True(t, ok)

		value, err := DecodeBytes(hexToBytes(t, "a1010000000000000000000000000000000000000000000000000000000000000000"))
		require.NoError(t, err)

		assert.Equal(t, Bytes, value.GetType())

		rawValue, ok := value.GetValue().([]byte)
		require.True(t, ok)

		assert.Equal(t, expectedNumber.Bytes(), big.NewInt(0).SetBytes(rawValue).Bytes())
	})
}

func TestDecode_EIP155(t *testing.T) {
	t.Parallel()

	// [nonce, gasPrice, gas, to, value, data, v, r, s] (9 elements)
	//
	// Nonce:    9
	// GasPrice: 20000000000
	// GasLimit: 21000
	// To:       0x3535353535353535353535353535353535353535
	// Value:    1000000000000000000
	// Data:     0x
	// V:        37
	// R:        0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276
	// S:        0x67cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83
	list, err := DecodeBytes(hexToBytes(t, "f86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"))
	require.NoError(t, err)

	require.Equal(t, List, list.GetType())

	listValues, _ := list.GetValue().([]Value)
	require.Len(t, listValues, 9)

	// Decode the nonce
	require.Equal(t, Bytes, listValues[0].GetType())

	nonceRaw, ok := listValues[0].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, nonceRaw)

	assert.Equal(t, uint64(9), big.NewInt(0).SetBytes(nonceRaw).Uint64())

	// Decode the gas price
	require.Equal(t, Bytes, listValues[1].GetType())

	gasPriceRaw, ok := listValues[1].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, gasPriceRaw)

	gasPrice := big.NewInt(0).SetBytes(gasPriceRaw)
	expectedGasPrice := big.NewInt(0).SetUint64(20000000000)

	assert.Zero(t, expectedGasPrice.Cmp(gasPrice))

	// Decode the gas limit
	require.Equal(t, Bytes, listValues[2].GetType())

	gasLimitRaw, ok := listValues[2].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, gasLimitRaw)

	gasLimit := big.NewInt(0).SetBytes(gasLimitRaw).Uint64()
	assert.Equal(t, uint64(21000), gasLimit)

	// Decode the recipient address
	require.Equal(t, Bytes, listValues[3].GetType())

	toRaw, ok := listValues[3].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, toRaw)

	assert.Equal(t, hexToBytes(t, "3535353535353535353535353535353535353535"), toRaw)

	// Decode the value
	require.Equal(t, Bytes, listValues[4].GetType())

	valueRaw, ok := listValues[4].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, valueRaw)

	value := big.NewInt(0).SetBytes(valueRaw)
	expectedValue, _ := big.NewInt(0).SetString("1000000000000000000", 10)

	assert.Zero(t, expectedValue.Cmp(value))

	// Decode the data
	require.Equal(t, Bytes, listValues[5].GetType())

	dataRaw, ok := listValues[5].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, dataRaw)

	assert.Empty(t, dataRaw)

	// Decode the v sig param
	require.Equal(t, Bytes, listValues[6].GetType())

	vRaw, ok := listValues[6].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, vRaw)

	v := big.NewInt(0).SetBytes(vRaw)
	expectedV := big.NewInt(0).SetUint64(37)

	assert.Zero(t, expectedV.Cmp(v))

	// Decode the r sig param
	require.Equal(t, Bytes, listValues[7].GetType())

	rRaw, ok := listValues[7].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, rRaw)

	r := big.NewInt(0).SetBytes(rRaw)
	expectedR, _ := big.NewInt(0).SetString("28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276", 16)

	assert.Zero(t, expectedR.Cmp(r))

	// Decode the s sig param
	require.Equal(t, Bytes, listValues[8].GetType())

	sRaw, ok := listValues[8].GetValue().([]byte)
	require.True(t, ok)
	require.NotNil(t, sRaw)

	s := big.NewInt(0).SetBytes(sRaw)
	expectedS, _ := big.NewInt(0).SetString("67cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83", 16)

	assert.Zero(t, expectedS.Cmp(s))
}
