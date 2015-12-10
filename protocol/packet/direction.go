package packet

type Direction int

const (
	DirectionServerbound Direction = iota
	DirectionClientbound
)
