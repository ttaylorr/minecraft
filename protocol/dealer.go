package protocol

import (
	"bytes"
	"errors"
	"reflect"
	"sync"

	"github.com/ttaylorr/minecraft/protocol/packet"
	"github.com/ttaylorr/minecraft/protocol/types"
	"github.com/ttaylorr/minecraft/util"
)

var (
	UnknownPacketError = errors.New("unknown packet type")
)

// A Dealer manages a set of `Rule`s and is able to decode and encode arbitrary
// data by finding and using applicable rules.
type Dealer struct {
	State State
	smu   sync.RWMutex
}

// NewDealer creates and returns a pointer to a new Dealer
func NewDealer() *Dealer {
	return &Dealer{}
}

func (d *Dealer) SetState(s State) {
	d.smu.Lock()
	defer d.smu.Unlock()

	d.State = s
}

// Decode decodes a packet coming from the client (sent to the server) into a
// packet "holder" type. (An example holder type is packet.Handshake). The
// types of the struct's fields are picked one by one (in order) and a field of
// data is decoded off of the stream and initialized into the corresponding
// field. If no matching decoder is found, an error is returned.
func (d *Dealer) Decode(p *packet.Packet) (v interface{}, err error) {
	htype := d.GetHolderType(p)
	if htype == nil {
		return nil, UnknownPacketError
	}

	inst := reflect.New(htype).Elem()

	data := bytes.NewBuffer(p.Data)

	for i := 0; i < inst.NumField(); i++ {
		f := inst.Field(i)

		typ, ok := f.Interface().(types.Type)
		if !ok {
			continue
		}

		v, err := typ.Decode(data)
		if err != nil {
			return nil, err
		}

		f.Set(reflect.ValueOf(v))
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
	out.Write(util.Uvarint(uint32(h.ID())))

	v := reflect.ValueOf(h)

	for i := 0; i < v.NumField(); i++ {
		ftype, ok := v.Field(i).Interface().(types.Type)
		if !ok {
			// XXX(taylor): special-casing
			ftype = types.JSON{V: v.Field(i).Interface()}
		}

		if _, err := ftype.Encode(out); err != nil {
			return nil, err
		}
	}

	return append(
		util.Uvarint(uint32(out.Len())),
		out.Bytes()...,
	), nil
}

// GetHolderType returns the `reflect.Type` associated with a particular packet
// sent to the server (from the client).
func (d *Dealer) GetHolderType(p *packet.Packet) reflect.Type {
	d.smu.RLock()
	defer d.smu.RUnlock()

	return GetPacket(p.Direction, d.State, p.ID)
}
