package packet

type Handshake struct {
	ProtocolVersion uint32
	ServerAddress   string
	ServerPort      uint16
	NextState       uint32
}

func (h Handshake) ID() int {
	return 0x00
}
