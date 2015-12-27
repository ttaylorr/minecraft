package types

import "io"

type UShort uint16

func (u UShort) Decode(r io.Reader) (interface{}, error) {
	buf := u.Buffer()
	if _, err := r.Read(buf); err != nil {
		return nil, err
	}

	return UShort(ByteOrder.Uint16(buf)), nil
}

func (u UShort) Encode(w io.Writer) (int, error) {
	buf := u.Buffer()
	ByteOrder.PutUint16(buf, uint16(u))

	return w.Write(buf)
}

func (u UShort) Buffer() []byte {
	return make([]byte, 2)
}
