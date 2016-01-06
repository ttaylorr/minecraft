package auth

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

func HashSession(serverId string, secret, pub []byte) string {
	sha1 := sha1.New()
	sha1.Write([]byte(serverId))
	sha1.Write(secret)
	sha1.Write(pub)

	hash := sha1.Sum(nil)

	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}
