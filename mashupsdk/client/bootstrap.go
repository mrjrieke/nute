package client

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"syscall"

	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"tini.com/nute/mashupsdk"
	sdk "tini.com/nute/mashupsdk"
)

var handshakeCredentials *sdk.MashupCredentials
var serverCredentials *sdk.MashupCredentials

var mashupContext *sdk.MashupContext

var handshakeCompleteChan chan bool

type MashupHandshakeServer struct {
	sdk.UnimplementedMashupServerServer
}

func (mhs *MashupHandshakeServer) Shake(ctx context.Context, in *sdk.MashupCredentials) (*sdk.MashupCredentials, error) {
	serverCredentials = in
	mashupCertBytes, err := mashupsdk.MashupKey.ReadFile("tls/mashup.crt")
	if err != nil {
		log.Fatalf("Couldn't load key: %v", err)
	}

	mashupClientCert, err := x509.ParseCertificate(mashupCertBytes)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Connect to the server for purposes of mashup api calls.
	mashupCertPool := x509.NewCertPool()
	mashupCertPool.AddCert(mashupClientCert)

	// Connect to it.

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(serverCredentials.Port)), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "", RootCAs: mashupCertPool})))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// Contact the server and print out its response.
	// User's of this library will benefit in following way:
	// 1. If current application shuts down, mashup
	// will also be told to shut down through Shutdown() api
	// call before this app exits.
	mashupContext.Client = sdk.NewMashupServerClient(conn)

	signalProcessor(mashupContext)

	handshakeCompleteChan <- true

	return &sdk.MashupCredentials{}, nil
}

// mashupInit -- starts mashup
func mashupInit(mashupGoodies map[string]interface{}) error {
	// exPath string, envParams []string, params []string

	var procAttr = syscall.ProcAttr{
		Dir:   ".",
		Env:   []string{"DISPLAY=:0.0"},
		Files: nil,
		Sys: &syscall.SysProcAttr{
			Setsid:     true,
			Foreground: false,
		},
	}
	var pid, err = syscall.ForkExec(mashupGoodies["MASHUP_PATH"].(string), mashupGoodies["PARAMS"].([]string), &procAttr)
	fmt.Println("Spawned proc", pid, err)
	mashupGoodies["PID"] = pid

	return err
}

func signalProcessor(mshCtx *sdk.MashupContext) {
	exitSignals := make(chan os.Signal, 4)
	signal.Notify(exitSignals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func(m *sdk.MashupContext) {
		<-exitSignals
		_, err := m.Client.Shutdown(m, &sdk.MashupEmpty{})
		if err != nil {
			log.Fatalf("Client shutdown failure: %v", err)
		}

	}(mshCtx)
}

func initContext(mashupGoodies map[string]interface{}) *sdk.MashupContext {
	mashupContext = &sdk.MashupContext{Context: context.Background(), MashupGoodies: mashupGoodies}
	// Initialize local server.
	mashupCertBytes, err := mashupsdk.MashupCert.ReadFile("tls/mashup.crt")
	if err != nil {
		log.Fatalf("Couldn't load cert: %v", err)
	}

	mashupKeyBytes, err := mashupsdk.MashupKey.ReadFile("tls/mashup.key")
	if err != nil {
		log.Fatalf("Couldn't load key: %v", err)
	}

	serverCert, err := tls.X509KeyPair(mashupCertBytes, mashupKeyBytes)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	creds := credentials.NewServerTLSFromCert(&serverCert)

	handshakeServer := grpc.NewServer(grpc.Creds(creds))
	lis, err := net.Listen("tcp", "localhost:0")
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	randomSha256 := sha256.Sum256(data)
	handshakeCredentials = &sdk.MashupCredentials{
		CallerToken: string(hex.EncodeToString([]byte(randomSha256[:]))),
		Port:        int64(lis.Addr().(*net.TCPAddr).Port),
	}

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	go handshakeServer.Serve(lis)

	jsonHandshakeCredentials, err := json.Marshal(handshakeCredentials)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}
	// Setup the onetime use handshake token...
	mashupGoodies["PARAMS"] = append(mashupGoodies["PARAMS"].([]string), "CREDS="+string(jsonHandshakeCredentials))

	// Start mashup..
	err = mashupInit(mashupGoodies)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}

	<-handshakeCompleteChan

	return mashupContext
}

func BootstrapInit(mashupPath string, envParams []string, params []string) *sdk.MashupContext {

	mashupGoodies := map[string]interface{}{}
	mashupGoodies["MASHUP_PATH"] = mashupPath
	mashupGoodies["ENV"] = envParams
	mashupGoodies["PARAMS"] = params

	return initContext(mashupGoodies)
}
