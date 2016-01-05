package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"

	"github.com/ttaylorr/minecraft/player"
	"github.com/ttaylorr/minecraft/protocol"
	mc "github.com/ttaylorr/minecraft/protocol/packet"
)

var (
	ClientDisconnectedError error = errors.New("client unexpectedly closed")
	VerifyKeyLength         int   = 16
)

type Authenticator struct {
	privateKey *rsa.PrivateKey
}

func New(pk *rsa.PrivateKey) *Authenticator {
	return &Authenticator{
		privateKey: pk,
	}
}

func (a *Authenticator) Login(c *protocol.Connection) (*player.Player, error) {
	verify, err := a.GenerateVerifyKey(VerifyKeyLength)
	if err != nil {
		return nil, err
	}

	for {
		packet, err := c.Next()
		if err != nil {
			return nil, err
		}

		switch v := packet.(type) {
		case mc.LoginStart:
			if err := a.loginStart(c, v, verify); err != nil {
				return nil, err
			}
		case mc.LoginEncryptionResponse:
			if err := a.onResponse(c, verify, v); err != nil {
				return nil, err
			}
			break
		}
	}

	return nil, nil
}

func (a *Authenticator) loginStart(c *protocol.Connection, ls mc.LoginStart, verify []byte) error {
	pub, err := x509.MarshalPKIXPublicKey(&a.privateKey.PublicKey)
	if err != nil {
		return err
	}

	req := mc.LoginEncryptionRequest{
		ServerID:  "",
		PublicKey: pub,
		VerifyKey: verify,
	}

	_, err = c.Write(req)
	return err
}

func (a *Authenticator) onResponse(c *protocol.Connection, verify []byte, r mc.LoginEncryptionResponse) error {
	//sharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, a.privateKey, r.SharedSecret)
	//if err != nil {
	//	return err
	//}

	tokenResp, err := rsa.DecryptPKCS1v15(rand.Reader, a.privateKey, r.VerifyToken)
	if err != nil {
		return err
	}

	if !bytes.Equal(tokenResp, verify) {
		return errors.New("verification tokens not equal")
	}

	return nil
}

func (a *Authenticator) GenerateVerifyKey(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)

	return b, err
}
