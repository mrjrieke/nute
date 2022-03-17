package server

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"log"
	"math/rand"
	"net"
	"strconv"

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

var serverCredentials *sdk.MashupCredentials

func (s *MashupServer) Shutdown(ctx context.Context, in *sdk.MashupEmpty) (*sdk.MashupEmpty, error) {
	os.Exit(-1)
	return &sdk.MashupEmpty{}, nil
}

func InitServer(creds string, insecure bool) {
	// Perform handshake...
	clientCredentials := sdk.MashupCredentials{}
	err := json.Unmarshal([]byte(creds), &clientCredentials)
	if err != nil {
		log.Fatalf("Malformed credentials: %s %v", creds, err)
	}

	go func(cc *sdk.MashupCredentials) {
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

		serverCredentials = &sdk.MashupCredentials{
			AuthToken:   cc.CallerToken,                                      // client's token.
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
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(cc.Port)), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "", RootCAs: mashupCertPool, InsecureSkipVerify: insecure})))
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

		_, err = mashupContext.Client.Shake(mashupContext.Context, serverCredentials)
		if err != nil {
			log.Fatalf("handshake failure: %v", err)
		}

		s.Serve(lis)

		sdk.RegisterMashupServerServer(s, &MashupServer{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}(&clientCredentials)
}
