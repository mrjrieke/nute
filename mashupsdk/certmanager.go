package mashupsdk

import (
	context "context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"math/rand"
)

//go:embed tls/mashup.crt
var MashupCert embed.FS

//go:embed tls/mashup.key
var MashupKey embed.FS

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
