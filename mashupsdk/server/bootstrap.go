package server

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"tini.com/nute/mashupsdk"
	sdk "tini.com/nute/mashupsdk"
)

// server is used to implement server.MashupServer.
type MashupServer struct {
	sdk.UnimplementedMashupServerServer
}

var clientConnectionConfigs *sdk.MashupConnectionConfigs
var serverConnectionConfigs *sdk.MashupConnectionConfigs

func (s *MashupServer) Shutdown(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupEmpty, error) {
	log.Printf("Shutdown called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")

	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		os.Exit(-1)
	}()

	log.Printf("Shutdown started")
	return &sdk.MashupEmpty{}, nil
}

func (s *MashupServer) OnResize(ctx context.Context, in *sdk.MashupDisplayBundle) (*sdk.MashupDisplayHint, error) {
	log.Printf("OnResize called")
	if in.GetAuthToken() != serverConnectionConfigs.AuthToken {
		return nil, errors.New("Auth failure")
	}
	// TODO: check credentials.
	//	in.MashupCredentials.AuthToken

	return nil, nil
}

func InitServer(creds string, insecure bool) {
	// Perform handshake...
	clientConnectionConfigs = &sdk.MashupConnectionConfigs{}
	err := json.Unmarshal([]byte(creds), clientConnectionConfigs)
	if err != nil {
		log.Fatalf("Malformed credentials: %s %v", creds, err)
	}
	log.Printf("Startup with insecure: %t\n", insecure)

	go func() {
		mashupCertBytes, err := mashupsdk.MashupCert.ReadFile("tls/mashup.crt")
		if err != nil {
			log.Fatalf("Couldn't load cert: %v", err)
		}

		mashupKeyBytes, err := mashupsdk.MashupKey.ReadFile("tls/mashup.key")
		if err != nil {
			log.Fatalf("Couldn't load key: %v", err)
		}

		cert, err := tls.X509KeyPair(mashupCertBytes, mashupKeyBytes)
		if err != nil {
			log.Fatalf("Couldn't construct key pair: %v", err)
		}
		creds := credentials.NewServerTLSFromCert(&cert)

		s := grpc.NewServer(grpc.Creds(creds))
		lis, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		data := make([]byte, 10)
		for i := range data {
			data[i] = byte(rand.Intn(256))
		}
		randomSha256 := sha256.Sum256(data)

		serverConnectionConfigs = &sdk.MashupConnectionConfigs{
			AuthToken:   clientConnectionConfigs.CallerToken,                 // client's token.
			CallerToken: string(hex.EncodeToString([]byte(randomSha256[:]))), // server token.
			Port:        int64(lis.Addr().(*net.TCPAddr).Port),
		}

		// Connect to the server for purposes of mashup api calls.
		mashupCertPool := x509.NewCertPool()
		mashupBlock, _ := pem.Decode([]byte(mashupCertBytes))
		mashupClientCert, err := x509.ParseCertificate(mashupBlock.Bytes)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		mashupCertPool.AddCert(mashupClientCert)

		// Send credentials back to client....
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(clientConnectionConfigs.Port)), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "", RootCAs: mashupCertPool, InsecureSkipVerify: insecure})))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		mashupContext := &sdk.MashupContext{Context: context.Background(), MashupGoodies: nil}

		// Contact the server and print out its response.
		// User's of this library will benefit in following way:
		// 1. If current application shuts down, mashup
		// will also be told to shut down through Shutdown() api
		// call before this app exits.
		mashupContext.Client = sdk.NewMashupServerClient(conn)

		go func() {
			// Async service initiation.
			log.Printf("Registering server.\n")

			sdk.RegisterMashupServerServer(s, &MashupServer{})

			log.Printf("Starting service.\n")
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()

		log.Printf("Handshake initiated.\n")

		_, err = mashupContext.Client.Shake(mashupContext.Context, serverConnectionConfigs)
		if err != nil {
			log.Printf("handshake failure: %v\n", err)
			panic(err)
		}
		log.Printf("Handshake complete.\n")

	}()
}
