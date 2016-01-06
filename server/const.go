package server

import (
	"crypto/rand"
	"crypto/rsa"
)

const (
	// PrivateKeySize determines the size of each server's generated private
	// key, in bytes. The Minecraft protocol mandates each be 1024 bytes, so
	// here we are!
	PrivateKeySize int = 1024
)

// Returns a new RSA private-key of length `PrivateKeySize`, precomputed.
func NewPrivateKey() (*rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, PrivateKeySize)
	if err != nil {
		return nil, err
	}
	priv.Precompute()

	return priv, nil
}
