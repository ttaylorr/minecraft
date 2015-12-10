package protocol

type Direction int

const (
	DirectionServerbound Direction = iota
	DirectionClientbound
)
