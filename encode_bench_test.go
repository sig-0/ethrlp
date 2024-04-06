package ethrlp

import "testing"

func benchmarkStringCommon(b *testing.B, value string) {
	b.Helper()

	for i := 0; i < b.N; i++ {
		_ = EncodeString(value)
	}
}

func BenchmarkEncode_String_Short(b *testing.B) {
	value := "Lorem ipsum dolor sit amet, consectetur adipisicing eli"

	b.ResetTimer()
	benchmarkStringCommon(b, value)
}

func BenchmarkEncode_String_Medium(b *testing.B) {
	value := "Lorem ipsum dolor sit amet, consectetur adipisicing elit"

	b.ResetTimer()
	benchmarkStringCommon(b, value)
}

func BenchmarkEncode_String_Long(b *testing.B) {
	value := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat"

	b.ResetTimer()
	benchmarkStringCommon(b, value)
}

func BenchmarkEncode_Array_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = EncodeArray([][]byte{
			EncodeString("aaa"),
			EncodeString("bbb"),
		})
	}
}

func BenchmarkEncode_Array_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = EncodeArray([][]byte{
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
	}
}

func BenchmarkEncode_Array_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = EncodeArray([][]byte{
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
	}
}

func BenchmarkEncode_Array_Nested_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = EncodeArray([][]byte{
			EncodeString("zw"),
			EncodeArray([][]byte{
				EncodeInt(4),
			}),
			EncodeInt(1),
		})
	}
}

func BenchmarkEncode_Array_Nested_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = EncodeArray([][]byte{
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
	}
}
