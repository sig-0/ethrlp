package ethrlp

type Type byte

func (t Type) String() string {
	switch t {
	case Bytes:
		return "Bytes"
	default:
		return "List"
	}
}

var (
	Bytes Type = 0x1
	List  Type = 0x2
)

// Value is a decoded data value
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
