package protocol

import (
	"bytes"
	"reflect"

	"github.com/ttaylorr/minecraft/protocol/packet"
	"github.com/ttaylorr/minecraft/protocol/rule"
	"github.com/ttaylorr/minecraft/util"
)

// A Dealer manages a set of `Rule`s and is able to decode and encode arbitrary
// data by finding and using applicable rules.
type Dealer struct {
	Rules []rule.Rule
}

// NewDealer creates and returns a pointer to a new Dealer, initialized with
// all of the `Rule`s passed to it.
func NewDealer(rules ...rule.Rule) *Dealer {
	return &Dealer{Rules: rules}
}

// DefaultDealer creates and returns a pointer to a new Dealer, initalized
// with all default and available `Rule`s.
func DefaultDealer() *Dealer {
	return NewDealer(
		rule.StringRule{},
		rule.UshortRule{},
		rule.UvarintRule{},
		rule.VarintRule{},
	)
}

// Decode decodes a packet coming from the client (sent to the server) into a
// packet "holder" type. (An example holder type is packet.Handshake). The
// types of the struct's fields are picked one by one (in order) and a field of
// data is decoded off of the stream and initialized into the corresponding
// field. If no matching decoder is found, an error is returned.
func (d *Dealer) Decode(p *packet.Packet) (v interface{}, err error) {
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

// Encode encodes a packet to be sent from the server to client into an array of
// bytes. The values of each field are encoded into the byte array according to
// their type.
//
// Once all of the data is serialized, the length and ID will be written in
// front of it, according to the specification given in `connection.go` from
// L41-L46.
//
// If an error occurs during encoding, an empty byte array and that
// error will be returned.
func (d *Dealer) Encode(h packet.Holder) ([]byte, error) {
	out := new(bytes.Buffer)
	out.Write(util.Uvarint(uint64(h.ID())))

	v := reflect.ValueOf(h)

	for i := 0; i < v.NumField(); i++ {
		fval := v.Field(i)

		rule := d.GetRule(fval.Type())
		encoded, err := rule.Encode(fval.Interface())
		if err != nil {
			return nil, err
		}

		out.Write(encoded)
	}

	return append(
		util.Uvarint(uint64(out.Len())),
		out.Bytes()...,
	), nil
}

// GetRule finds the first matching rule given a particular type. It queries the
// `Rule#AppliesTo` method and returns the first matching one. If no matching
// `Rule`s are found, a value of `nil` is returned instead.
func (d *Dealer) GetRule(typ reflect.Type) rule.Rule {
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
func (d *Dealer) GetHolderType(p *packet.Packet) reflect.Type {
	return packet.Packets[p.ID]
}
