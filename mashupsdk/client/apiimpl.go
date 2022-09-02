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

	"github.com/mrjrieke/nute/mashupsdk"
	sdk "github.com/mrjrieke/nute/mashupsdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type MashupClient struct {
	sdk.UnimplementedMashupServerServer
	mashupApiHandler mashupsdk.MashupApiHandler
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
		os.Exit(-1)
	}()

	log.Println("Shutdown initiated.")
	return &sdk.MashupEmpty{}, nil
}

// Shake - Implementation of the handshake.  During the callback from
// the mashup, construct new more permanent set of credentials to be shared.
func (c *MashupClient) Shake(ctx context.Context, in *sdk.MashupConnectionConfigs) (*sdk.MashupConnectionConfigs, error) {
	log.Printf("Shake called")
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
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(serverConnectionConfigs.Port)), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "", RootCAs: mashupCertPool, InsecureSkipVerify: *insecure})))
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
	log.Printf("Handshake complete.\n")

	return clientConnectionConfigs, nil
}

// Shutdown -- handles request to shut down the mashup.
func (c *MashupClient) GetMashupElements(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupDetailedElementBundle, error) {
	log.Printf("GetMashupElements called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		return c.mashupApiHandler.GetMashupElements()
	} else {
		log.Printf("No api handler provided.")
	}
	return nil, nil
}

func (c *MashupClient) UpsertMashupElementsState(ctx context.Context, in *sdk.MashupElementStateBundle) (*sdk.MashupElementStateBundle, error) {
	log.Printf("UpsertMashupElementsState called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		return c.mashupApiHandler.UpsertMashupElementsState(in)
	} else {
		log.Printf("No api handler provided.")
	}
	return nil, nil
}

func (c *MashupClient) UpsertMashupElements(ctx context.Context, in *sdk.MashupDetailedElementBundle) (*sdk.MashupDetailedElementBundle, error) {
	log.Printf("UpsertMashupElements called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		return c.mashupApiHandler.UpsertMashupElements(in)
	} else {
		log.Printf("No api handler provided.")
	}
	return nil, nil
}

func (c *MashupClient) OnResize(ctx context.Context, in *sdk.MashupDisplayBundle) (*sdk.MashupDisplayHint, error) {
	log.Printf("OnResize called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	displayHint := in.MashupDisplayHint
	log.Printf("Received resize: %d %d %d %d\n", displayHint.Xpos, displayHint.Ypos, displayHint.Width, displayHint.Height)
	if c.mashupApiHandler != nil {
		log.Printf("Delegate to api handler.")
		c.mashupApiHandler.OnResize(displayHint)
	} else {
		log.Printf("No api handler provided.")
	}
	return displayHint, nil
}
