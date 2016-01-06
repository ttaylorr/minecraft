package handshake_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/protocol/packet"
	"github.com/ttaylorr/minecraft/protocol/types"
	"github.com/ttaylorr/minecraft/server/handshake"
)

func TestNew(t *testing.T) {
	h := handshake.New()
	assert.IsType(t, h, &handshake.Handshaker{})
}

func TestPongHasMatchingPayload(t *testing.T) {
	h := handshake.New()
	payload := types.Long(rand.Int63())
	ping := packet.StatusPing{payload}

	pong := h.GetPong(ping)

	assert.Equal(t, pong.Payload, ping.Payload)
}

func TestStatusMatchesMOTD(t *testing.T) {
	// TODO(taylor)
	t.Skip()
}
