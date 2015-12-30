package protocol

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/ttaylorr/minecraft/protocol/packet"
)

// Connection represents a uni-directional connection from client to server.
type Connection struct {
	d  *Dealer
	rw io.ReadWriter
}

// NewConnection serves as the builder function for type Connection. It takes in
// a reader which, when read from, yeilds data sent by the "client".
func NewConnection(rw io.ReadWriter) *Connection {
	return &Connection{d: NewDealer(), rw: rw}
}

func (c *Connection) SetState(state State) {
	c.d.SetState(state)
}

func (c *Connection) Next() (interface{}, error) {
	p, err := c.packet()
	if err != nil {
		return nil, err
	}

	return c.d.Decode(p)
}

func (c *Connection) Write(h packet.Holder) (int, error) {
	data, err := c.d.Encode(h)
	if err != nil {
		return -1, nil
	}

	return c.rw.Write(data)
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
	r := bufio.NewReader(c.rw)

	size, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, size)
	_, err = io.ReadAtLeast(r, buffer, int(size))
	if err != nil {
		return nil, err
	}

	id, offset := binary.Uvarint(buffer)

	return &packet.Packet{
		ID:        int(id),
		Direction: packet.DirectionServerbound,
		Data:      buffer[offset:],
	}, nil
}
