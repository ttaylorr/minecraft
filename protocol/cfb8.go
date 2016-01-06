/*
   Copyright 2013 Matthew Collins (purggames@gmail.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package protocol

import (
	"crypto/cipher"
)

/*
	Allow for AES streams
*/
type cfb8 struct {
	c                cipher.Block
	blockSize        int
	iv, iv_real, tmp []byte
	de               bool
}

func newCFB8(c cipher.Block, iv []byte, decrypt bool) *cfb8 {
	if len(iv) != 16 {
		panic("bad iv length!")
	}
	cp := make([]byte, 256)
	copy(cp, iv)
	return &cfb8{
		c:         c,
		blockSize: c.BlockSize(),
		iv:        cp[:16],
		iv_real:   cp,
		tmp:       make([]byte, 16),
		de:        decrypt,
	}
}

func newCFB8Decrypt(c cipher.Block, iv []byte) *cfb8 {
	return newCFB8(c, iv, true)
}

func newCFB8Encrypt(c cipher.Block, iv []byte) *cfb8 {
	return newCFB8(c, iv, false)
}

func (cf *cfb8) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		val := src[i]
		cf.c.Encrypt(cf.tmp, cf.iv)
		val = val ^ cf.tmp[0]

		if cap(cf.iv) >= 17 {
			cf.iv = cf.iv[1:17]
		} else {
			copy(cf.iv_real, cf.iv[1:])
			cf.iv = cf.iv_real[:16]
		}

		if cf.de {
			cf.iv[15] = src[i]
		} else {
			cf.iv[15] = val
		}
		dst[i] = val
	}
}
