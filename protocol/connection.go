package protocol

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"

	"github.com/ttaylorr/minecraft/protocol/packet"
)

// Connection represents a uni-directional connection from client to server.
type Connection struct {
	d *Dealer

	in  io.Reader
	out io.Writer
}

// NewConnection serves as the builder function for type Connection. It takes in
// a reader which, when read from, yeilds data sent by the "client".
func NewConnection(rw io.ReadWriter) *Connection {
	return &Connection{d: NewDealer(), in: rw, out: rw}
}

// SetState changes the protocol state (see https://wiki.vg) of the connection
// between the Server and Client. StateChagnes are proxied down into the Dealer,
// and happen in a sync fashion.
func (c *Connection) SetState(state State) {
	c.d.SetState(state)
}

// Next will read and return the Go struct representing the data contained in
// the next packet.
func (c *Connection) Next() (packet.Holder, error) {
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

	return c.out.Write(data)
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
// ... (TODO)
//
//
// If an error is experienced in reading the packet from the io.Reader `r`, then
// a nil pointer will be returned and the error will be propogated up.
func (c *Connection) packet() (*packet.Packet, error) {
	r := bufio.NewReader(c.in)

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

// Encrypt begins encrypting the connection with an AES/CFB8 cipher as according
// to the most-current Minecraft protocol. An AES cipher is seeded from the
// de-encrypted shared-secret token, and the initial vector of both the reader
// and writer are set to that same shared-secret.
func (c *Connection) Encrypt(secret []byte) error {
	aes, err := aes.NewCipher(secret)
	if err != nil {
		return err
	}

	c.in = cipher.StreamReader{
		R: c.in,
		S: newCFB8Decrypt(aes, secret),
	}

	c.out = cipher.StreamWriter{
		W: c.out,
		S: newCFB8Encrypt(aes, secret),
	}

	return nil
}
