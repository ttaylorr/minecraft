package mcio

import (
	"bytes"
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/mcio/types"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

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

func GetFieldType(field reflect.StructField) types.Type {
	return types.GetType(FieldType(field))
}

func FieldType(field reflect.StructField) string {
	return field.Tag.Get("type")
}
