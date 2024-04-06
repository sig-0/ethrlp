package ethrlp

import (
	"testing"
)

func benchmarkDecodeCommon(b *testing.B, input []byte) {
	for i := 0; i < b.N; i++ {
		_, _ = DecodeBytes(input)
	}
}

func BenchmarkDecode_String_Short(b *testing.B) {
	value := hexToBytes(b, "B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974")

	b.ResetTimer()
	benchmarkDecodeCommon(b, value)
}

func BenchmarkDecode_String_Medium(b *testing.B) {
	value := hexToBytes(b, "B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974")

	b.ResetTimer()
	benchmarkDecodeCommon(b, value)
}

func BenchmarkDecode_String_Long(b *testing.B) {
	value := hexToBytes(b, "B904004C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E73656374657475722061646970697363696E6720656C69742E20437572616269747572206D6175726973206D61676E612C20737573636970697420736564207665686963756C61206E6F6E2C20696163756C697320666175636962757320746F72746F722E2050726F696E20737573636970697420756C74726963696573206D616C6573756164612E204475697320746F72746F7220656C69742C2064696374756D2071756973207472697374697175652065752C20756C7472696365732061742072697375732E204D6F72626920612065737420696D70657264696574206D6920756C6C616D636F7270657220616C6971756574207375736369706974206E6563206C6F72656D2E2041656E65616E2071756973206C656F206D6F6C6C69732C2076756C70757461746520656C6974207661726975732C20636F6E73657175617420656E696D2E204E756C6C6120756C74726963657320747572706973206A7573746F2C20657420706F73756572652075726E6120636F6E7365637465747572206E65632E2050726F696E206E6F6E20636F6E76616C6C6973206D657475732E20446F6E65632074656D706F7220697073756D20696E206D617572697320636F6E67756520736F6C6C696369747564696E2E20566573746962756C756D20616E746520697073756D207072696D697320696E206661756369627573206F726369206C756374757320657420756C74726963657320706F737565726520637562696C69612043757261653B2053757370656E646973736520636F6E76616C6C69732073656D2076656C206D617373612066617563696275732C2065676574206C6163696E6961206C616375732074656D706F722E204E756C6C61207175697320756C747269636965732070757275732E2050726F696E20617563746F722072686F6E637573206E69626820636F6E64696D656E74756D206D6F6C6C69732E20416C697175616D20636F6E73657175617420656E696D206174206D65747573206C75637475732C206120656C656966656E6420707572757320656765737461732E20437572616269747572206174206E696268206D657475732E204E616D20626962656E64756D2C206E6571756520617420617563746F72207472697374697175652C206C6F72656D206C696265726F20616C697175657420617263752C206E6F6E20696E74657264756D2074656C6C7573206C65637475732073697420616D65742065726F732E20437261732072686F6E6375732C206D65747573206163206F726E617265206375727375732C20646F6C6F72206A7573746F20756C747269636573206D657475732C20617420756C6C616D636F7270657220766F6C7574706174")

	b.ResetTimer()
	benchmarkDecodeCommon(b, value)
}

func BenchmarkDecode_Array_Small(b *testing.B) {
	encoding := EncodeArray([][]byte{
		EncodeString("aaa"),
		EncodeString("bbb"),
	})

	b.ResetTimer()
	benchmarkDecodeCommon(b, encoding)
}

func BenchmarkDecode_Array_Medium(b *testing.B) {
	encoding := EncodeArray([][]byte{
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
	})

	b.ResetTimer()
	benchmarkDecodeCommon(b, encoding)
}

func BenchmarkDecode_Array_Large(b *testing.B) {
	encoding := EncodeArray([][]byte{
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
	})

	b.ResetTimer()
	benchmarkDecodeCommon(b, encoding)
}

func BenchmarkDecode_Array_Nested_Short(b *testing.B) {
	encoding := EncodeArray([][]byte{
		EncodeString("zw"),
		EncodeArray([][]byte{
			EncodeInt(4),
		}),
		EncodeInt(1),
	})

	b.ResetTimer()
	benchmarkDecodeCommon(b, encoding)
}

func BenchmarkDecode_Array_Nested_Long(b *testing.B) {
	encoding := EncodeArray([][]byte{
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
	})

	b.ResetTimer()
	benchmarkDecodeCommon(b, encoding)
}
