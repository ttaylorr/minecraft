package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"io"

	"github.com/ttaylorr/minecraft/player"
	"github.com/ttaylorr/minecraft/protocol"
	"github.com/ttaylorr/minecraft/protocol/packet"
)

var (
	ClientDisconnectedError error = errors.New("client unexpectedly closed")
)

type Authenticator struct {
	privateKey *rsa.PrivateKey
}

func New(pk *rsa.PrivateKey) {
	return &Authenticator{
		privateKey: pk,
	}
}

func (a *Authenticator) Login(c *protocol.Connection) (*player.Player, error) {
	for {
		packet, err := c.Next()
		if err == io.EOF {
			return nil, ClientDisconnectedError
		}

		if err = a.handlePacket(packet); err != nil {
			return nil, err
		}
	}
}

func (a *Authenticator) handlePacket(p *packet.Holder) error {
	verify, err := a.GenerateVerifyKey()
	if err != nil {
		return err
	}

	switch v := p.(type) {
	case packets.Handshake:
		a.c.SetState(protocol.State(uint8(v.State)))
	case packets.LoginStart:
		pub, err := x509.MarshalPKIXPublicKey(&a.privateKey.PublicKey)
		if err != nil {
			return err
		}

		req := packets.LoginEncryptionRequest{
			ServerID:  "",
			PublicKey: pub,
			VerifyKey: verify,
		}
		c.Write(req)
	case packets.LoginEncryptionResponse:
		sharedSecret, err := rsa.DecryptPKC1v15(rand.Reader, a.privateKey, v.SharedSecret)
		if err != nil {
			return err
		}

		tokenResp, err := rsa.DecryptPKC1v15(rand.Reader, a.privateKey, v.VerifyToken)
		if err != nil {
			return err
		}

		if !bytes.Equal(tokenResp, verify) {
			return errors.New("verification tokens not equal")
		}
	}
}

func (a *Authenticator) GenerateVerifyKey(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)

	return b, err
}
