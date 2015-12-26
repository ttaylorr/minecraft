package packet

type Handshake struct {
	ProtocolVersion uint64
	ServerAddress   string
	ServerPort      uint16
	NextState       uint64
}
