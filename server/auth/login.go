package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/ttaylorr/minecraft/player"
	"github.com/ttaylorr/minecraft/protocol"
	mc "github.com/ttaylorr/minecraft/protocol/packet"
	"github.com/ttaylorr/minecraft/protocol/types"
)

var (
	ClientDisconnectedError error = errors.New("client unexpectedly closed")
	VerifyKeyLength         int   = 16
)

type Authenticator struct {
	Yggdrasil *Yggdrasil

	privateKey *rsa.PrivateKey
}

func New(pk *rsa.PrivateKey) (*Authenticator, error) {
	a := &Authenticator{
		privateKey: pk,
	}

	pub, err := a.PublicKey()
	if err != nil {
		return nil, err
	}
	a.Yggdrasil = NewYggdrasil(pub)

	return a, nil
}

func (a *Authenticator) Login(c *protocol.Connection) (*player.Player, error) {
	var username string

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
			username = string(v.Username)
			if err := a.loginStart(c, v, verify); err != nil {
				return nil, err
			}
		case mc.LoginEncryptionResponse:
			if err := a.onResponse(c, username, verify, v); err != nil {
				return nil, err
			}
			break
		}
	}

	return nil, nil
}

func (a *Authenticator) PublicKey() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&a.privateKey.PublicKey)
}

func (a *Authenticator) loginStart(c *protocol.Connection, ls mc.LoginStart,
	verify []byte) error {

	pub, err := a.PublicKey()
	if err != nil {
		return err
	}

	_, err = c.Write(mc.LoginEncryptionRequest{
		ServerID:  "",
		PublicKey: pub,
		VerifyKey: verify,
	})

	return err
}

func (a *Authenticator) onResponse(c *protocol.Connection, username string,
	verify []byte, r mc.LoginEncryptionResponse) error {

	secret, _ := a.decrypt(r.SharedSecret)
	check, _ := a.decrypt(r.VerifyToken)

	if !bytes.Equal(verify, check) {
		return errors.New("verification token not equal")
	}

	session, err := a.Yggdrasil.GetSession(username, secret)

	uuid := session.ID
	uuid = fmt.Sprintf("%s-%s-%s-%s-%s", uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:32])

	c.Encrypt(secret)

	if _, err = c.Write(mc.LoginSuccess{
		UUID:     types.String(uuid),
		Username: types.String(username),
	}); err != nil {
		return err
	}

	return nil
}

func (a *Authenticator) decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, a.privateKey, data)
}

func (a *Authenticator) GenerateVerifyKey(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)

	return b, err
}
