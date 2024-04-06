package ethrlp

import (
	"encoding/hex"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testingT is a sugar interface for sharing
// helpers between *testing.T and *testing.B
type testingT interface {
	Helper()

	require.TestingT
	assert.TestingT
}

// hexToBytes converts a string hex representation to a bytes array
// (without the leading 0x)
func hexToBytes(t testingT, input string) []byte {
	t.Helper()

	r := strings.NewReplacer("\t", "", " ", "", "\n", "")

	data, err := hex.DecodeString(r.Replace(input))
	require.NoError(t, err)

	return data
}

func TestEncode_Byte(t *testing.T) {
	t.Parallel()

	t.Run("Byte in [0x00, 0x7f] range", func(t *testing.T) {
		t.Parallel()

		for i := 0x00; i <= 0x7f; i++ {
			assert.Equal(t, []byte{byte(i)}, EncodeByte(byte(i)))
		}
	})

	t.Run("Byte in the [0x80, 0xff] range", func(t *testing.T) {
		t.Parallel()

		for i := 0x80; i <= 0xff; i++ {
			assert.Equal(t, []byte{0x81, byte(i)}, EncodeByte(byte(i)))
		}
	})
}

func TestEncode_Bytes(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name           string
		input          []byte
		expectedOutput []byte
	}{
		{
			"Empty bytes",
			nil,
			hexToBytes(t, "80"),
		},
		{
			"Input bytes <=55B in length",
			[]byte("hello world"),
			hexToBytes(t, "8B68656C6C6F20776F726C64"),
		},
		{
			"Input bytes >55B in length",
			[]byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit"),
			hexToBytes(t, "B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974"),
		},
		{
			"Input bytes >55B in length, long",
			[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat"),
			hexToBytes(t, "B904004C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E73656374657475722061646970697363696E6720656C69742E20437572616269747572206D6175726973206D61676E612C20737573636970697420736564207665686963756C61206E6F6E2C20696163756C697320666175636962757320746F72746F722E2050726F696E20737573636970697420756C74726963696573206D616C6573756164612E204475697320746F72746F7220656C69742C2064696374756D2071756973207472697374697175652065752C20756C7472696365732061742072697375732E204D6F72626920612065737420696D70657264696574206D6920756C6C616D636F7270657220616C6971756574207375736369706974206E6563206C6F72656D2E2041656E65616E2071756973206C656F206D6F6C6C69732C2076756C70757461746520656C6974207661726975732C20636F6E73657175617420656E696D2E204E756C6C6120756C74726963657320747572706973206A7573746F2C20657420706F73756572652075726E6120636F6E7365637465747572206E65632E2050726F696E206E6F6E20636F6E76616C6C6973206D657475732E20446F6E65632074656D706F7220697073756D20696E206D617572697320636F6E67756520736F6C6C696369747564696E2E20566573746962756C756D20616E746520697073756D207072696D697320696E206661756369627573206F726369206C756374757320657420756C74726963657320706F737565726520637562696C69612043757261653B2053757370656E646973736520636F6E76616C6C69732073656D2076656C206D617373612066617563696275732C2065676574206C6163696E6961206C616375732074656D706F722E204E756C6C61207175697320756C747269636965732070757275732E2050726F696E20617563746F722072686F6E637573206E69626820636F6E64696D656E74756D206D6F6C6C69732E20416C697175616D20636F6E73657175617420656E696D206174206D65747573206C75637475732C206120656C656966656E6420707572757320656765737461732E20437572616269747572206174206E696268206D657475732E204E616D20626962656E64756D2C206E6571756520617420617563746F72207472697374697175652C206C6F72656D206C696265726F20616C697175657420617263752C206E6F6E20696E74657264756D2074656C6C7573206C65637475732073697420616D65742065726F732E20437261732072686F6E6375732C206D65747573206163206F726E617265206375727375732C20646F6C6F72206A7573746F20756C747269636573206D657475732C20617420756C6C616D636F7270657220766F6C7574706174"),
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expectedOutput, EncodeBytes(testCase.input))
		})
	}
}

func TestEncode_Bool(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name           string
		expectedOutput []byte
		input          bool
	}{
		{
			"true",
			hexToBytes(t, "01"),
			true,
		},
		{
			"false",
			hexToBytes(t, "80"),
			false,
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expectedOutput, EncodeBool(testCase.input))
		})
	}
}

func TestEncode_Int(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name           string
		expectedOutput []byte
		input          int64
	}{
		{
			"0",
			hexToBytes(t, "80"),
			0,
		},
		{
			"1",
			hexToBytes(t, "01"),
			1,
		},
		{
			"127",
			hexToBytes(t, "7F"),
			127,
		},
		{
			"128",
			hexToBytes(t, "8180"),
			128,
		},
		{
			"256",
			hexToBytes(t, "820100"),
			256,
		},
		{
			"1024",
			hexToBytes(t, "820400"),
			1024,
		},
		{
			"0xFFFFFF",
			hexToBytes(t, "83FFFFFF"),
			0xFFFFFF,
		},
		{
			"0xFFFFFFFF",
			hexToBytes(t, "84FFFFFFFF"),
			0xFFFFFFFF,
		},
		{
			"0xFFFFFFFFFF",
			hexToBytes(t, "85FFFFFFFFFF"),
			0xFFFFFFFFFF,
		},
		{
			"0xFFFFFFFFFFFF",
			hexToBytes(t, "86FFFFFFFFFFFF"),
			0xFFFFFFFFFFFF,
		},
		{
			"0xFFFFFFFFFFFFFF",
			hexToBytes(t, "87FFFFFFFFFFFFFF"),
			0xFFFFFFFFFFFFFF,
		},
		{
			"max int64",
			hexToBytes(t, "887FFFFFFFFFFFFFFF"),
			math.MaxInt64,
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expectedOutput, EncodeInt(testCase.input))
		})
	}
}

func TestEncode_Uint(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name           string
		expectedOutput []byte
		input          uint64
	}{
		{
			"0",
			hexToBytes(t, "80"),
			0,
		},
		{
			"1",
			hexToBytes(t, "01"),
			1,
		},
		{
			"127",
			hexToBytes(t, "7F"),
			127,
		},
		{
			"128",
			hexToBytes(t, "8180"),
			128,
		},
		{
			"256",
			hexToBytes(t, "820100"),
			256,
		},
		{
			"1024",
			hexToBytes(t, "820400"),
			1024,
		},
		{
			"0xFFFFFF",
			hexToBytes(t, "83FFFFFF"),
			0xFFFFFF,
		},
		{
			"0xFFFFFFFF",
			hexToBytes(t, "84FFFFFFFF"),
			0xFFFFFFFF,
		},
		{
			"0xFFFFFFFFFF",
			hexToBytes(t, "85FFFFFFFFFF"),
			0xFFFFFFFFFF,
		},
		{
			"0xFFFFFFFFFFFF",
			hexToBytes(t, "86FFFFFFFFFFFF"),
			0xFFFFFFFFFFFF,
		},
		{
			"0xFFFFFFFFFFFFFF",
			hexToBytes(t, "87FFFFFFFFFFFFFF"),
			0xFFFFFFFFFFFFFF,
		},
		{
			"max uint64",
			hexToBytes(t, "88FFFFFFFFFFFFFFFF"),
			math.MaxUint64,
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expectedOutput, EncodeUint(testCase.input))
		})
	}
}

func TestEncode_BigInt(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name           string
		input          string
		expectedOutput []byte
	}{
		{
			"0x102030405060708090A0B0C0D0E0F2",
			"83729609699884896815286331701780722",
			hexToBytes(t, "8F102030405060708090A0B0C0D0E0F2"),
		},
		{
			"0x0100020003000400050006000700080009000A000B000C000D000E01",
			"105315505618206987246253880190783558935785933862974822347068935681",
			hexToBytes(t, "9C0100020003000400050006000700080009000A000B000C000D000E01"),
		},
		{
			"0x010000000000000000000000000000000000000000000000000000000000000000",
			"115792089237316195423570985008687907853269984665640564039457584007913129639936",
			hexToBytes(t, "A1010000000000000000000000000000000000000000000000000000000000000000"),
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			number, ok := big.NewInt(0).SetString(testCase.input, 10)
			require.True(t, ok)

			assert.Equal(t, testCase.expectedOutput, EncodeBigInt(number))
		})
	}
}

func TestEncode_String(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name           string
		input          string
		expectedOutput []byte
	}{
		{
			"empty string",
			"",
			hexToBytes(t, "80"),
		},
		{
			"string < 7F",
			"\x7E",
			hexToBytes(t, "7E"),
		},
		{
			"string == 7F",
			"\x7F",
			hexToBytes(t, "7F"),
		},
		{
			"string > 7F",
			"\x80",
			hexToBytes(t, "8180"),
		},
		{
			"simple word",
			"dog",
			hexToBytes(t, "83646F67"),
		},
		{
			"short text",
			"Lorem ipsum dolor sit amet, consectetur adipisicing eli",
			hexToBytes(t, "B74C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C69"),
		},
		{
			"medium text",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			hexToBytes(t, "B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974"),
		},
		{
			"long text",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat",
			hexToBytes(t, "B904004C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E73656374657475722061646970697363696E6720656C69742E20437572616269747572206D6175726973206D61676E612C20737573636970697420736564207665686963756C61206E6F6E2C20696163756C697320666175636962757320746F72746F722E2050726F696E20737573636970697420756C74726963696573206D616C6573756164612E204475697320746F72746F7220656C69742C2064696374756D2071756973207472697374697175652065752C20756C7472696365732061742072697375732E204D6F72626920612065737420696D70657264696574206D6920756C6C616D636F7270657220616C6971756574207375736369706974206E6563206C6F72656D2E2041656E65616E2071756973206C656F206D6F6C6C69732C2076756C70757461746520656C6974207661726975732C20636F6E73657175617420656E696D2E204E756C6C6120756C74726963657320747572706973206A7573746F2C20657420706F73756572652075726E6120636F6E7365637465747572206E65632E2050726F696E206E6F6E20636F6E76616C6C6973206D657475732E20446F6E65632074656D706F7220697073756D20696E206D617572697320636F6E67756520736F6C6C696369747564696E2E20566573746962756C756D20616E746520697073756D207072696D697320696E206661756369627573206F726369206C756374757320657420756C74726963657320706F737565726520637562696C69612043757261653B2053757370656E646973736520636F6E76616C6C69732073656D2076656C206D617373612066617563696275732C2065676574206C6163696E6961206C616375732074656D706F722E204E756C6C61207175697320756C747269636965732070757275732E2050726F696E20617563746F722072686F6E637573206E69626820636F6E64696D656E74756D206D6F6C6C69732E20416C697175616D20636F6E73657175617420656E696D206174206D65747573206C75637475732C206120656C656966656E6420707572757320656765737461732E20437572616269747572206174206E696268206D657475732E204E616D20626962656E64756D2C206E6571756520617420617563746F72207472697374697175652C206C6F72656D206C696265726F20616C697175657420617263752C206E6F6E20696E74657264756D2074656C6C7573206C65637475732073697420616D65742065726F732E20437261732072686F6E6375732C206D65747573206163206F726E617265206375727375732C20646F6C6F72206A7573746F20756C747269636573206D657475732C20617420756C6C616D636F7270657220766F6C7574706174"),
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expectedOutput, EncodeString(testCase.input))
		})
	}
}

func TestEncode_Array(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name           string
		input          [][]byte
		expectedOutput []byte
	}{
		{
			"Empty array",
			nil,
			hexToBytes(t, "c0"),
		},
		{
			name: "Populated array <=55B in length",
			input: [][]byte{
				EncodeString("hello"),
				EncodeString("world"),
			},
			expectedOutput: hexToBytes(t, "CC8568656C6C6F85776F726C64"),
		},
		{
			name: "Populated array >55B in length",
			input: [][]byte{
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
			expectedOutput: hexToBytes(t, "F83C836161618362626283636363836464648365656583666666836767678368686883696969836A6A6A836B6B6B836C6C6C836D6D6D836E6E6E836F6F6F"),
		},
	}

	for _, testCase := range testTable {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expectedOutput, EncodeArray(testCase.input))
		})
	}
}

// Test cases copied over from the official Ethereum RLP testing suite:
// https://github.com/ethereum/tests/blob/develop/RLPTests/rlptest.json
func TestEncode_EthereumTests(t *testing.T) {
	t.Parallel()

	t.Run("emptystring", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, hexToBytes(t, "80"), EncodeString(""))
	})

	t.Run("bytestring00", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, hexToBytes(t, "00"), EncodeBytes([]byte{0x0}))
	})

	t.Run("bytestring01", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, hexToBytes(t, "01"), EncodeBytes([]byte{0x1}))
	})

	t.Run("bytestring7F", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, hexToBytes(t, "7f"), EncodeByte(0x7F))
	})

	t.Run("shortstring", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, hexToBytes(t, "83646f67"), EncodeString("dog"))
	})

	t.Run("shortstring2", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"b74c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c69",
			),
			EncodeString("Lorem ipsum dolor sit amet, consectetur adipisicing eli"),
		)
	})

	t.Run("longstring", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"b8384c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c6974",
			),
			EncodeString("Lorem ipsum dolor sit amet, consectetur adipisicing elit"),
		)
	})

	t.Run("longstring2", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"b904004c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e73656374657475722061646970697363696e6720656c69742e20437572616269747572206d6175726973206d61676e612c20737573636970697420736564207665686963756c61206e6f6e2c20696163756c697320666175636962757320746f72746f722e2050726f696e20737573636970697420756c74726963696573206d616c6573756164612e204475697320746f72746f7220656c69742c2064696374756d2071756973207472697374697175652065752c20756c7472696365732061742072697375732e204d6f72626920612065737420696d70657264696574206d6920756c6c616d636f7270657220616c6971756574207375736369706974206e6563206c6f72656d2e2041656e65616e2071756973206c656f206d6f6c6c69732c2076756c70757461746520656c6974207661726975732c20636f6e73657175617420656e696d2e204e756c6c6120756c74726963657320747572706973206a7573746f2c20657420706f73756572652075726e6120636f6e7365637465747572206e65632e2050726f696e206e6f6e20636f6e76616c6c6973206d657475732e20446f6e65632074656d706f7220697073756d20696e206d617572697320636f6e67756520736f6c6c696369747564696e2e20566573746962756c756d20616e746520697073756d207072696d697320696e206661756369627573206f726369206c756374757320657420756c74726963657320706f737565726520637562696c69612043757261653b2053757370656e646973736520636f6e76616c6c69732073656d2076656c206d617373612066617563696275732c2065676574206c6163696e6961206c616375732074656d706f722e204e756c6c61207175697320756c747269636965732070757275732e2050726f696e20617563746f722072686f6e637573206e69626820636f6e64696d656e74756d206d6f6c6c69732e20416c697175616d20636f6e73657175617420656e696d206174206d65747573206c75637475732c206120656c656966656e6420707572757320656765737461732e20437572616269747572206174206e696268206d657475732e204e616d20626962656e64756d2c206e6571756520617420617563746f72207472697374697175652c206c6f72656d206c696265726f20616c697175657420617263752c206e6f6e20696e74657264756d2074656c6c7573206c65637475732073697420616d65742065726f732e20437261732072686f6e6375732c206d65747573206163206f726e617265206375727375732c20646f6c6f72206a7573746f20756c747269636573206d657475732c20617420756c6c616d636f7270657220766f6c7574706174",
			),
			EncodeString("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat"),
		)
	})

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"80",
			),
			EncodeUint(0),
		)
	})

	t.Run("smallint", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"01",
			),
			EncodeInt(1),
		)
	})

	t.Run("smallint2", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"10",
			),
			EncodeInt(16),
		)
	})

	t.Run("smallint3", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"4f",
			),
			EncodeInt(79),
		)
	})

	t.Run("smallint4", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"7f",
			),
			EncodeInt(127),
		)
	})

	t.Run("mediumint1", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"8180",
			),
			EncodeInt(128),
		)
	})

	t.Run("mediumint2", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"8203e8",
			),
			EncodeInt(1000),
		)
	})

	t.Run("mediumint3", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"830186a0",
			),
			EncodeInt(100000),
		)
	})

	t.Run("mediumint4", func(t *testing.T) {
		t.Parallel()

		number, ok := big.NewInt(0).SetString("83729609699884896815286331701780722", 10)
		require.True(t, ok)

		assert.Equal(
			t,
			hexToBytes(
				t,
				"8f102030405060708090a0b0c0d0e0f2",
			),
			EncodeBigInt(number),
		)
	})

	t.Run("mediumint5", func(t *testing.T) {
		t.Parallel()

		number, ok := big.NewInt(0).SetString(
			"105315505618206987246253880190783558935785933862974822347068935681",
			10,
		)
		require.True(t, ok)

		assert.Equal(
			t,
			hexToBytes(
				t,
				"9c0100020003000400050006000700080009000a000b000c000d000e01",
			),
			EncodeBigInt(number),
		)
	})

	t.Run("emptylist", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"c0",
			),
			EncodeArray([][]byte{}),
		)
	})

	t.Run("stringlist", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"cc83646f6783676f6483636174",
			),
			EncodeArray([][]byte{
				EncodeString("dog"),
				EncodeString("god"),
				EncodeString("cat"),
			}),
		)
	})

	t.Run("multilist", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"c6827a77c10401",
			),
			EncodeArray([][]byte{
				EncodeString("zw"),
				EncodeArray([][]byte{
					EncodeInt(4),
				}),
				EncodeInt(1),
			}),
		)
	})

	t.Run("shortListMax1", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"f784617364668471776572847a78637684617364668471776572847a78637684617364668471776572847a78637684617364668471776572",
			),
			EncodeArray([][]byte{
				EncodeString("asdf"),
				EncodeString("qwer"),
				EncodeString("zxcv"),
				EncodeString("asdf"),
				EncodeString("qwer"),
				EncodeString("zxcv"),
				EncodeString("asdf"),
				EncodeString("qwer"),
				EncodeString("zxcv"),
				EncodeString("asdf"),
				EncodeString("qwer"),
			}),
		)
	})

	t.Run("longList1", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"f840cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376",
			),
			EncodeArray([][]byte{
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
			}),
		)
	})

	t.Run("longList2", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"f90200cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376cf84617364668471776572847a786376",
			),
			EncodeArray([][]byte{
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
				EncodeArray([][]byte{
					EncodeString("asdf"),
					EncodeString("qwer"),
					EncodeString("zxcv"),
				}),
			}),
		)
	})

	t.Run("listsoflists", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"c4c2c0c0c0",
			),
			// [ [ [], [] ], [] ]
			EncodeArray([][]byte{
				EncodeArray([][]byte{
					EncodeArray([][]byte{}),
					EncodeArray([][]byte{}),
				}),
				EncodeArray([][]byte{}),
			}),
		)
	})

	t.Run("listsoflists2", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"c7c0c1c0c3c0c1c0",
			),
			// [ [], [[]], [ [], [[]] ] ]
			EncodeArray([][]byte{
				EncodeArray([][]byte{}),
				EncodeArray([][]byte{
					EncodeArray([][]byte{}),
				}),
				EncodeArray([][]byte{
					EncodeArray([][]byte{}),
					EncodeArray([][]byte{
						EncodeArray([][]byte{}),
					}),
				}),
			}),
		)
	})

	t.Run("dictTest1", func(t *testing.T) {
		t.Parallel()

		assert.Equal(
			t,
			hexToBytes(
				t,
				"ecca846b6579318476616c31ca846b6579328476616c32ca846b6579338476616c33ca846b6579348476616c34",
			),
			EncodeArray([][]byte{
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
			}),
		)
	})

	t.Run("bigint", func(t *testing.T) {
		t.Parallel()

		number, ok := big.NewInt(0).SetString(
			"115792089237316195423570985008687907853269984665640564039457584007913129639936",
			10,
		)
		require.True(t, ok)

		assert.Equal(
			t,
			hexToBytes(
				t,
				"a1010000000000000000000000000000000000000000000000000000000000000000",
			),
			EncodeBigInt(number),
		)
	})
}
