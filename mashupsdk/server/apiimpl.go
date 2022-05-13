package server

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"tini.com/nute/mashupsdk"
	sdk "tini.com/nute/mashupsdk"
)

// server is used to implement server.MashupServer.
type MashupServer struct {
	sdk.UnimplementedMashupServerServer
	mashupApiHandler mashupsdk.MashupApiHandler
}

func GetClientAuthToken() string {
	if clientConnectionConfigs != nil {
		return clientConnectionConfigs.AuthToken
	} else {
		return ""
	}
}

func GetServerAuthToken() string {
	if serverConnectionConfigs != nil {
		return serverConnectionConfigs.AuthToken
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

func (s *MashupServer) UpsertMashupElements(ctx context.Context, in *sdk.MashupDetailedElementBundle) (*sdk.MashupElementStateBundle, error) {
	log.Printf("UpsertMashupElements called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if s.mashupApiHandler != nil {
		log.Printf("UpsertMashupElements Delegate to api handler.")
		return s.mashupApiHandler.UpsertMashupElements(in)
	}
	return nil, status.Errorf(codes.Unimplemented, "method UpsertMashupElements not implemented")
}

func (s *MashupServer) UpsertMashupElementsState(ctx context.Context, in *sdk.MashupElementStateBundle) (*sdk.MashupElementStateBundle, error) {
	log.Printf("UpsertMashupElementsState called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if s.mashupApiHandler != nil {
		log.Printf("UpsertMashupElementsState Delegate to api handler.")
		return s.mashupApiHandler.UpsertMashupElementsState(in)
	}
	return nil, nil
}
