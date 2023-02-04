package server

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/mrjrieke/nute/mashupsdk"
	sdk "github.com/mrjrieke/nute/mashupsdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	log.Println("Shutdown complete.")
	return &sdk.MashupEmpty{}, nil
}

func (s *MashupServer) ResetStates(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupEmpty, error) {
	log.Println("ResetStates called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if s.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		s.mashupApiHandler.ResetStates()
	}

	log.Println("ResetStates complete.")
	return &sdk.MashupEmpty{}, nil
}

// OnDisplayChange -- handles a request from the client to resize.
func (s *MashupServer) OnDisplayChange(ctx context.Context, in *sdk.MashupDisplayBundle) (*sdk.MashupDisplayHint, error) {
	log.Printf("OnDisplayChange called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		log.Printf("OnDisplayChange auth failure.")
		return nil, errors.New("Auth failure")
	}
	displayHint := in.MashupDisplayHint
	log.Printf("Received resize: %d %d %d %d\n", displayHint.Xpos, displayHint.Ypos, displayHint.Width, displayHint.Height)

	if s.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		s.mashupApiHandler.OnDisplayChange(displayHint)
	}

	return nil, nil
}

func (s *MashupServer) GetElements(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupDetailedElementBundle, error) {
	log.Printf("GetElements called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if s.mashupApiHandler != nil {
		log.Printf("GetElements Delegate to api handler.")
		return s.mashupApiHandler.GetElements()
	}
	return nil, status.Errorf(codes.Unimplemented, "method GetMashupElements not implemented")
}

func (s *MashupServer) UpsertElements(ctx context.Context, in *sdk.MashupDetailedElementBundle) (*sdk.MashupDetailedElementBundle, error) {
	log.Printf("UpsertElements called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if s.mashupApiHandler != nil {
		log.Printf("UpsertElements Delegate to api handler.")
		return s.mashupApiHandler.UpsertElements(in)
	}
	return nil, status.Errorf(codes.Unimplemented, "method UpsertElements not implemented")
}

func (s *MashupServer) TweakStates(ctx context.Context, in *sdk.MashupElementStateBundle) (*sdk.MashupElementStateBundle, error) {
	log.Printf("TweakStates called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		log.Printf("TweakStates Auth failure.")
		return nil, errors.New("Auth failure")
	}
	if s.mashupApiHandler != nil {
		log.Printf("TweakStates Delegate to api handler.")
		return s.mashupApiHandler.TweakStates(in)
	}
	return nil, nil
}
