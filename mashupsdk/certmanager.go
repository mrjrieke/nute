package mashupsdk

import (
	context "context"
	"embed"
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
