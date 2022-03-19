package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	sdk "tini.com/nute/mashupsdk"
)

func initSignalProcessor(mshCtx *sdk.MashupContext) {
	exitSignals := make(chan os.Signal, 4)
	signal.Notify(exitSignals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func(m *sdk.MashupContext) {
		<-exitSignals
		// TODO: Send real credentials.
		_, err := m.Client.Shutdown(m, &sdk.MashupEmpty{AuthToken: GetServerAuthToken()})
		if err != nil {
			log.Fatalf("Possibly using self signed cert.  Consider using -insecure flag for developing.  Client shutdown failure: %v", err)
		}
		log.Printf("Client shutting down.")
		os.Exit(0)
	}(mshCtx)
}
