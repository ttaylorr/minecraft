package protocol

import (
	"bytes"
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/packet"
)

type Rule interface {
	AppliesTo(typ reflect.Type) bool
	Decode(r *bytes.Buffer) (interface{}, error)
}

type Decoder struct {
	Rules []Rule
}

func NewDecoder(rules ...Rule) *Decoder {
	return &Decoder{Rules: rules}
}

func DefaultDecoder() *Decoder {
	return NewDecoder(
		StringRule{},
		UshortRule{},
		UvarintRule{},
		VarintRule{},
	)
}

func (d *Decoder) Decode(p *packet.Packet) (v interface{}, err error) {
	typ := d.GetHolderType(p)
	inst := reflect.New(typ).Elem()

	data := bytes.NewBuffer(p.Data)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		rule := d.GetRule(field.Type)

		decoded, err := rule.Decode(data)
		if err != nil {
			return nil, err
		}

		inst.Field(i).Set(reflect.ValueOf(decoded))
	}

	return inst.Interface(), nil
}

func (d *Decoder) GetRule(typ reflect.Type) Rule {
	for _, rule := range d.Rules {
		if !rule.AppliesTo(typ) {
			continue
		}

		return rule
	}

	return nil
}

func (d *Decoder) GetHolderType(p *packet.Packet) reflect.Type {
	return packet.Packets[p.ID]
}
