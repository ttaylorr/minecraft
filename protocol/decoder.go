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

// A Decoder manages a set of `Rule`s and is able to decode arbitrary data by
// finding and using applicable rules.
type Decoder struct {
	Rules []Rule
}

// NewDecoder creates and returns a pointer to a new Decoder, initialized with
// all of the `Rule`s passed to it.
func NewDecoder(rules ...Rule) *Decoder {
	return &Decoder{Rules: rules}
}

// DefaultDecoder creates and returns a pointer to a new Decoder, initalized
// with all default and available `Rule`s.
func DefaultDecoder() *Decoder {
	return NewDecoder(
		StringRule{},
		UshortRule{},
		UvarintRule{},
		VarintRule{},
	)
}

// Decode decodes a packet coming from the client (sent to the server) into a
// packet "holder" type. (An example holder type is packet.Handshake). The
// types of the struct's fields are picked one by one (in order) and a field of
// data is decoded off of the stream and initialized into the corresponding
// field. If no matching decoder is found, an error is returned.
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

// GetRule finds the first matching rule given a particular type. It queries the
// `Rule#AppliesTo` method and returns the first matching one. If no matching
// `Rule`s are found, a value of `nil` is returned instead.
func (d *Decoder) GetRule(typ reflect.Type) Rule {
	for _, rule := range d.Rules {
		if !rule.AppliesTo(typ) {
			continue
		}

		return rule
	}

	return nil
}

// GetHolderType returns the `reflect.Type` associated with a particular packet
// sent to the server (from the client).
func (d *Decoder) GetHolderType(p *packet.Packet) reflect.Type {
	return packet.Packets[p.ID]
}
