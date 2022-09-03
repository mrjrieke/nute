//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"embed"
	"flag"
	"log"
	"os"

	"github.com/mrjrieke/nute/custos/custosworld"
	"github.com/mrjrieke/nute/g3nd/data"
	"github.com/mrjrieke/nute/mashupsdk"
)

var worldCompleteChan chan bool

//go:embed tls/mashup.crt
var mashupCert embed.FS

//go:embed tls/mashup.key
var mashupKey embed.FS

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	headless := flag.Bool("headless", false, "Run headless")
	flag.Parse()
	worldLog, err := os.OpenFile("custos.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	detailedElements := data.GetExampleLibrary()

	custosWorld := custosworld.NewCustosWorldApp(*headless, detailedElements, nil)

	custosWorld.InitServer(*callerCreds, *insecure)

	// Initialize the main window.
	custosWorld.InitMainWindow()

	<-worldCompleteChan
}
