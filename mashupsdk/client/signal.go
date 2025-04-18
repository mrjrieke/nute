package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	sdk "github.com/mrjrieke/nute-core/mashupsdk"
)

func initSignalProcessor(mshCtx *sdk.MashupContext) {
	exitSignals := make(chan os.Signal, 4)
	signal.Notify(exitSignals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func(m *sdk.MashupContext) {
		<-exitSignals
		_, err := m.Client.Shutdown(m, &sdk.MashupEmpty{AuthToken: GetServerAuthToken()})
		if err != nil {
			log.Fatalf("Possibly using self signed cert.  Consider using -tls-skip-validation flag for developing.  Client shutdown failure: %v", err)
		}
		log.Printf("Client shutting down.")
		os.Exit(0)
	}(mshCtx)
}
