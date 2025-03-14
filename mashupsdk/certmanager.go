package mashupsdk

import (
	context "context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"math/rand"
)

var MashupCertBytes []byte

var MashupKeyBytes []byte

func InitCertKeyPair(mc embed.FS, mk embed.FS) error {
	var err error
	MashupCertBytes, err = mc.ReadFile("tls/mashup.crt")
	if err != nil {
		return err
	}

	MashupKeyBytes, err = mk.ReadFile("tls/mashup.key")
	if err != nil {
		return err
	}
	return nil
}

func InitCertKeyPairBytes(mashupCertBytes []byte, mashupKeyBytes []byte) {
	MashupCertBytes = mashupCertBytes
	MashupKeyBytes = mashupKeyBytes
}

type MashupContext struct {
	context.Context
	Client        MashupServerClient
	MashupGoodies map[string]interface{}
}

func GenAuthToken() string {

	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	randomSha256 := sha256.Sum256(data)

	return string(hex.EncodeToString([]byte(randomSha256[:])))
}
