package mcio

import (
	"bytes"
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/mcio/types"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(p *packet.Packet) (interface{}, error) {
	typ, val := d.NewPacket(p.ID)
	r := bytes.NewBuffer(p.Data)

	for i := 0; i < typ.NumField(); i++ {
		if err := d.DecodeField(typ, val, i, r); err != nil {
			return nil, err
		}
	}

	return val.Interface(), nil
}

func (d *Decoder) NewPacket(id int) (reflect.Type, reflect.Value) {
	typ := packet.Packets[id]
	return typ, reflect.New(typ).Elem()
}

func (d *Decoder) DecodeField(typ reflect.Type, val reflect.Value, field int, r *bytes.Buffer) error {
	ftype := d.GetFieldType(typ.Field(field))
	decoded, err := ftype.Decode(r)
	if err != nil {
		return err
	}

	d.SetFieldValue(val, field, decoded)

	return nil
}

func (d *Decoder) SetFieldValue(holder reflect.Value, field int, value interface{}) {
	v := reflect.ValueOf(value)
	holder.Field(field).Set(v)
}

func (d *Decoder) GetFieldType(field reflect.StructField) types.Type {
	return types.GetType(d.FieldType(field))
}

func (d *Decoder) FieldType(field reflect.StructField) string {
	return field.Tag.Get("type")
}
