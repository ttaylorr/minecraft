package protocol

import (
	"bytes"
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/packet"
)

// A Rule represents a particular decoding method, usually bound to a type. It
// gives information regarding how to decode ~~(and encode)~~ a particular piece
// of data ~~to and~~ from a `*bytes.Buffer`.
//
// Only after a call to the `AppliesTo` method returns true can it be said that
// the Rule is allowed to work on a given value.
type Rule interface {
	// AppliesTo predicates which type(s) a particular rule can work upon.
	//
	// Typically AppliesTo only passes for a single type, and it is usually
	// implemented as:
	//
	// ```go
	// func (r RuleImpl) AppliesTo(typ reflect.Type) bool {
	//     return typ.Kind == reflect.SomeKind
	// }
	AppliesTo(typ reflect.Type) bool

	// Decode reads from the given *bytes.Buffer and returns the decoded
	// contents of a single field of data. If an error is encountered during
	// read-time, or if the data is invalid post read-time, an error will be
	// thrown.
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
