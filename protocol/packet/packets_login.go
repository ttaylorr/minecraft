package packet

import "github.com/ttaylorr/minecraft/protocol/types"

type (
	LoginStart struct {
		Name types.String
	}

	LoginEncryptionRequest struct {
		ServerID  types.String
		PublicKey types.ByteArray
		VerifyKey types.ByteArray
	}

	LoginEncryptionResponse struct {
		SharedSecret types.ByteArray
		VerifyToken  types.ByteArray
	}

	LoginSuccess struct {
		UUID     types.String
		Username types.String
	}

	LoginSetCompression struct {
		Threshold types.UVarint
	}
)

func (_ LoginStart) ID() int              { return 0x01 }
func (_ LoginEncryptionRequest) ID() int  { return 0x01 }
func (_ LoginEncryptionResponse) ID() int { return 0x01 }
func (_ LoginSuccess) ID() int            { return 0x02 }
func (_ LoginSetCompression) ID() int     { return 0x03 }
