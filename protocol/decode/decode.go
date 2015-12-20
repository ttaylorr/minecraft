package decode

import (
	"bytes"
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/decode/types"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

var (
	Types = map[string]Type{
		"varint":  types.Varint{},
		"uvarint": types.Uvarint{},
		"string":  types.String{},
		"ushort":  types.Ushort{},
	}
)

type Type interface {
	Decode(r *bytes.Buffer) (interface{}, error)
	// Encode(v interface{}) []byte
}

func Decode(p *packet.Packet) (interface{}, error) {
	typ, val := NewPacket(p.ID)
	r := bytes.NewBuffer(p.Data)

	for i := 0; i < typ.NumField(); i++ {
		err := DecodeField(typ, val, i, r)
		if err != nil {
			return nil, err
		}
	}

	return val.Interface(), nil
}

func NewPacket(id int) (reflect.Type, reflect.Value) {
	typ := packet.Packets[id]
	return typ, reflect.New(typ).Elem()
}

func DecodeField(typ reflect.Type, val reflect.Value, field int, r *bytes.Buffer) error {
	ftype := GetFieldType(typ.Field(field))
	decoded, err := ftype.Decode(r)
	if err != nil {
		return err
	}

	SetFieldValue(val, field, decoded)

	return nil
}

func SetFieldValue(holder reflect.Value, field int, value interface{}) {
	v := reflect.ValueOf(value)
	holder.Field(field).Set(v)
}

func GetFieldType(field reflect.StructField) Type {
	return Types[FieldType(field)]
}

func FieldType(field reflect.StructField) string {
	return field.Tag.Get("type")
}
