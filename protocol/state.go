package protocol

// State represnts a protocol state over the network. With each changing state,
// the relative Packet IDs are changed.
type State uint8

const (
	HandshakeState State = iota
	StatusState
	LoginState
)
