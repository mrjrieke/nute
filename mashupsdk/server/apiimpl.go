package server

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	sdk "tini.com/nute/mashupsdk"
)

// server is used to implement server.MashupServer.
type MashupServer struct {
	sdk.UnimplementedMashupServerServer
	onResizeCb func(displayHint *sdk.MashupDisplayHint)
}

func GetClientAuthToken() string {
	if clientConnectionConfigs != nil {
		return clientConnectionConfigs.AuthToken
	} else {
		return ""
	}
}

// Shutdown -- handles request to shut down the mashup.
func (s *MashupServer) Shutdown(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupEmpty, error) {
	log.Println("Shutdown called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		os.Exit(-1)
	}()

	log.Println("Shutdown initiated.")
	return &sdk.MashupEmpty{}, nil
}

// OnResize -- handles a request from the client to resize.
func (s *MashupServer) OnResize(ctx context.Context, in *sdk.MashupDisplayBundle) (*sdk.MashupDisplayHint, error) {
	log.Printf("OnResize called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	displayHint := in.MashupDisplayHint
	log.Printf("Received resize: %d %d %d %d\n", displayHint.Xpos, displayHint.Ypos, displayHint.Width, displayHint.Height)

	s.onResizeCb(displayHint)

	return nil, nil
}
