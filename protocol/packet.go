package protocol

// Type Packet encapsulates a raw packet sent between the Minecraft server and
// the Minecraft client.
type Packet struct {
	// ID is the ID of the packet according to the latest version of the
	// Minecraft protocol as implemented by Mojang AB.
	ID int

	// Direction gives information regarding who sent the packet. Packets
	// sent from server to client will always be tagged DirectionClientbound
	// whereas messages passed from client ot server will be
	// DirectionServerbound.
	Direction Direction

	// Data contains the uncompressed array of bytes encoding the body of
	// the packet itself.
	Data []byte
}
