package ethrlp

type Type byte

var (
	Bytes Type = 0x1
	List  Type = 0x2
)

type Value interface {
	// GetType returns the type information of the value
	GetType() Type

	// GetValue returns the raw uncast value
	GetValue() any
}

type BytesValue struct {
	value []byte
}

func (b *BytesValue) GetType() Type {
	return Bytes
}

func (b *BytesValue) GetValue() any {
	return b.value
}

type ListValue struct {
	values []Value
}

func (a *ListValue) GetType() Type {
	return List
}

func (a *ListValue) GetValue() any {
	return a.values
}
