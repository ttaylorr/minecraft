package protocol

type State uint8

const (
	HandshakeState State = iota
	StatusState
	LoginState
)
