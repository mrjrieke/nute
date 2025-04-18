package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mrjrieke/nute-core/mashupsdk"
	sdk "github.com/mrjrieke/nute-core/mashupsdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

var clientDialOptions grpc.DialOption = grpc.EmptyDialOption{}

type MashupClient struct {
	sdk.UnimplementedMashupServerServer
	mashupApiHandler mashupsdk.MashupApiHandler
}

func InitDialOptions(dialOption grpc.DialOption) {
	clientDialOptions = dialOption
}

func GetServerAuthToken() string {
	if serverConnectionConfigs != nil {
		return serverConnectionConfigs.AuthToken
	} else {
		return ""
	}
}

// Shutdown -- handles request to shut down the mashup.
func (s *MashupClient) Shutdown(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupEmpty, error) {
	log.Println("Shutdown called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		log.Printf("Client shutting down.")
		os.Exit(-1)
	}()

	log.Println("Shutdown initiated.")
	return &sdk.MashupEmpty{}, nil
}

// CollaborateInit - Implementation of the handshake.  During the callback from
// the mashup, construct new more permanent set of credentials to be shared.
func (c *MashupClient) CollaborateInit(ctx context.Context, in *sdk.MashupConnectionConfigs) (*sdk.MashupConnectionConfigs, error) {
	log.Printf("CollaborateInit called")
	if in.GetAuthToken() != handshakeConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	serverConnectionConfigs = &sdk.MashupConnectionConfigs{
		AuthToken: in.CallerToken,
		Port:      in.Port,
	}

	if mashupCertBytes == nil {
		log.Fatalf("Cert not initialized.")
	}
	mashupBlock, _ := pem.Decode([]byte(mashupCertBytes))
	mashupClientCert, err := x509.ParseCertificate(mashupBlock.Bytes)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Connect to the server for purposes of mashup api calls.
	mashupCertPool := x509.NewCertPool()
	mashupCertPool.AddCert(mashupClientCert)

	log.Printf("Initiating connection to server with insecure: %t\n", *insecure)
	// Connect to it.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(serverConnectionConfigs.Port)), clientDialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "", RootCAs: mashupCertPool, InsecureSkipVerify: *insecure})))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("Connection to server established.\n")

	// Contact the server and print out its response.
	// User's of this library will benefit in following way:
	// 1. If current application shuts down, mashup
	// will also be told to shut down through Shutdown() api
	// call before this app exits.
	mashupContext.Client = sdk.NewMashupServerClient(conn)
	log.Printf("Initiate signal handler.\n")

	initSignalProcessor(mashupContext)
	log.Printf("Signal handler initialized.\n")

	go func() {
		handshakeCompleteChan <- true
	}()

	clientConnectionConfigs = &sdk.MashupConnectionConfigs{
		AuthToken: sdk.GenAuthToken(), // client token.
		Port:      handshakeConnectionConfigs.Port,
	}
	log.Printf("CollaborateInit complete.\n")

	return clientConnectionConfigs, nil
}

// Shutdown -- handles request to shut down the mashup.
func (c *MashupClient) GetElements(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupDetailedElementBundle, error) {
	log.Printf("GetElements called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		return c.mashupApiHandler.GetElements()
	} else {
		log.Printf("No api handler provided.")
	}
	return nil, nil
}

func (c *MashupClient) TweakStates(ctx context.Context, in *sdk.MashupElementStateBundle) (*sdk.MashupElementStateBundle, error) {
	log.Printf("TweakStates called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		log.Printf("Auth failure.")
		return nil, errors.New("Auth failure")
	}
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		return c.mashupApiHandler.TweakStates(in)
	} else {
		log.Printf("No api handler provided.")
	}
	return nil, nil
}

func (c *MashupClient) TweakStatesByMotiv(ctx context.Context, in *mashupsdk.Motiv) (*emptypb.Empty, error) {
	log.Printf("TweakStatesByMotiv called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		log.Printf("Auth failure.")
		return nil, errors.New("Auth failure")
	}
	if c.mashupApiHandler != nil {
		log.Printf("TweakStatesByMotiv Delegate to api handler.")
		return c.mashupApiHandler.TweakStatesByMotiv(in)
	} else {
		log.Printf("TweakStatesByMotiv No api handler provided.")
	}
	return nil, nil
}

func (c *MashupClient) UpsertElements(ctx context.Context, in *sdk.MashupDetailedElementBundle) (*sdk.MashupDetailedElementBundle, error) {
	log.Printf("UpsertElements called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		return c.mashupApiHandler.UpsertElements(in)
	} else {
		log.Printf("No api handler provided.")
	}
	return nil, nil
}

func (c *MashupClient) OnDisplayChange(ctx context.Context, in *sdk.MashupDisplayBundle) (*sdk.MashupDisplayHint, error) {
	log.Printf("OnDisplayChange called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		log.Printf("OnDisplayChange auth failure.")
		return nil, errors.New("Auth failure")
	}
	displayHint := in.MashupDisplayHint
	log.Printf("Received resize: %d %d %d %d\n", displayHint.Xpos, displayHint.Ypos, displayHint.Width, displayHint.Height)
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		c.mashupApiHandler.OnDisplayChange(displayHint)
	} else {
		log.Printf("No api handler provided.")
	}
	return displayHint, nil
}
