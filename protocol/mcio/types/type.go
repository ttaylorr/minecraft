package types

import "bytes"

var (
	types = map[string]Type{
		"varint":  Varint{},
		"uvarint": Uvarint{},
		"string":  String{},
		"ushort":  Ushort{},
	}
)

func GetType(id string) Type {
	return types[id]
}

type Type interface {
	Decode(r *bytes.Buffer) (interface{}, error)
	// Encode(v interface{}) []byte
}
