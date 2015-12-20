package packet

type Handshake struct {
	ProtocolVersion uint64 `type:"uvarint"`
	ServerAddress   string `type:"string"`
	ServerPort      uint16 `type:"ushort"`
	NextState       uint64 `type:"uvarint"`
}
