package server

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	sdk "tini.com/nute/mashupsdk"
)

// server is used to implement server.MashupServer.
type MashupServer struct {
	sdk.UnimplementedMashupServerServer
	mashupApiHandler MashupApiHandler
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

	if s.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		s.mashupApiHandler.OnResize(displayHint)
	}

	return nil, nil
}

func (mhs *MashupServer) UpsertMashupSociety(ctx context.Context, in *sdk.MashupSocietyBundle) (*sdk.MashupSocietyStateBundle, error) {
	log.Printf("UpsertMashupSociety called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	// TODO: Implement.

	return nil, status.Errorf(codes.Unimplemented, "method UpsertMashupSociety not implemented")
}

func (mhs *MashupServer) UpsertMashupSocietyState(ctx context.Context, in *sdk.MashupSocietyStateBundle) (*sdk.MashupSocietyStateBundle, error) {
	log.Printf("UpsertMashupSocietyState called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	// TODO: Implement.

	return nil, nil
}
