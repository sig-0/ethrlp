package ethrlp

import "fmt"

func ExampleEncodeBool() {
	fmt.Printf("%X\n", EncodeBool(true))
	fmt.Printf("%X\n", EncodeBool(false))

	// Output:
	// 01
	// 80
}

func ExampleEncodeInt() {
	fmt.Printf("%X\n", EncodeInt(0xFFFFFF))

	// Output:
	// 83FFFFFF
}

func ExampleEncodeUint() {
	fmt.Printf("%X\n", EncodeUint(0xFFFFFFFFFFFFFF))

	// Output:
	// 87FFFFFFFFFFFFFF
}

func ExampleEncodeString() {
	fmt.Printf("%X\n", EncodeString("hello world"))

	// Output:
	// 8B68656C6C6F20776F726C64
}

func ExampleEncodeByte() {
	fmt.Printf("%X\n", EncodeByte(0x5))
	fmt.Printf("%X\n", EncodeByte(0x80))

	// Output:
	// 05
	// 8180
}

func ExampleEncodeBytes() {
	fmt.Printf("%X\n", EncodeBytes([]byte("hello world")))
	fmt.Printf("%X\n", EncodeBytes([]byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit")))

	// Output:
	// 8B68656C6C6F20776F726C64
	// B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974
}

func ExampleEncodeArray() {
	fmt.Printf("%X\n", EncodeArray([][]byte{
		EncodeString("hello"),
		EncodeString("world"),
	}))

	fmt.Printf("%X\n", EncodeArray([][]byte{
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
	}))

	// Output:
	// CC8568656C6C6F85776F726C64
	// F83C836161618362626283636363836464648365656583666666836767678368686883696969836A6A6A836B6B6B836C6C6C836D6D6D836E6E6E836F6F6F
}
