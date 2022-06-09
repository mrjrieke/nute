//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"flag"
	"log"
	"os"

	"tini.com/nute/g3nd/g3nworld"
)

var worldCompleteChan chan bool

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("insecure", false, "Skip server validation")
	flag.Parse()
	worldLog, err := os.OpenFile("world.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(worldLog)

	worldApp := g3nworld.NewWorldApp()

	worldApp.InitServer(*callerCreds, *insecure)

	// Initialize the main window.
	go worldApp.InitMainWindow()

	<-worldCompleteChan
}
