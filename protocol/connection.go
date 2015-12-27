package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

	"github.com/ttaylorr/minecraft/protocol/packet"
)

var (
	ErrTooFewBytes = errors.New("too few bytes read")
)

// Connection represents a uni-directional connection from client to server.
type Connection struct {
	d *Dealer
	r io.Reader
}

// NewConnection serves as the builder function for type Connection. It takes in
// a reader which, when read from, yeilds data sent by the "client".
func NewConnection(r io.Reader) *Connection {
	return &Connection{d: DefaultDealer(), r: r}
}

func (c *Connection) Next() (interface{}, error) {
	p, err := c.packet()
	if err != nil {
		return nil, err
	}

	return c.d.Decode(p)
}

// Next reads and decodes the next Packet on the stream. Packets are expected to
// be in the following format (as described on
// http://wiki.vg/Protocol#Without_compression:
//
// Without compression:
//   | Field Name | Field Type | Field Notes                        |
//   | ---------- | ---------- | ---------------------------------- |
//   | Length     | Uvarint    | Represents length of <id> + <data> |
//   | ID         | Uvarint    |                                    |
//   | Data       | []byte     |                                    |
//
// With compression:
// ...
//
// If an error is experienced in reading the packet from the io.Reader `r`, then
// a nil pointer will be returned and the error will be propogated up.
func (c *Connection) packet() (*packet.Packet, error) {
	r := bufio.NewReader(c.r)

	size, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	// TODO(ttaylorr): extract this to a package `util`
	buffer := make([]byte, size)
	read, err := r.Read(buffer)
	if err != nil {
		return nil, err
	} else if read < int(size) {
		return nil, ErrTooFewBytes
	}

	id, offset := binary.Uvarint(buffer)

	return &packet.Packet{
		ID:        int(id),
		Direction: packet.DirectionServerbound,
		Data:      buffer[offset:],
	}, nil
}
