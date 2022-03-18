package client

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
var insecure *bool

var handshakeCompleteChan chan bool

type MashupHandshakeServer struct {
	sdk.UnimplementedMashupServerServer
}

var mashupCertBytes []byte

func (mhs *MashupHandshakeServer) Shake(ctx context.Context, in *sdk.MashupCredentials) (*sdk.MashupCredentials, error) {
	log.Printf("Handshake initiated.\n")
	serverCredentials = in

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
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(serverCredentials.Port)), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "", RootCAs: mashupCertPool, InsecureSkipVerify: *insecure})))
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

	signalProcessor(mashupContext)
	log.Printf("Signal handler initialized.\n")

	go func() { handshakeCompleteChan <- true }()
	log.Printf("Handshake complete.\n")

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
	params := []string{mashupGoodies["MASHUP_PATH"].(string)}
	params = append(params, mashupGoodies["PARAMS"].([]string)...)

	var pid, err = syscall.ForkExec(mashupGoodies["MASHUP_PATH"].(string), params, &procAttr)
	log.Println("Spawned proc", pid, err)
	mashupGoodies["PID"] = pid

	return err
}

func signalProcessor(mshCtx *sdk.MashupContext) {
	exitSignals := make(chan os.Signal, 4)
	signal.Notify(exitSignals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func(m *sdk.MashupContext) {
		<-exitSignals
		// TODO: Send real credentials.
		_, err := m.Client.Shutdown(m, &sdk.MashupCredentials{})
		if err != nil {
			log.Fatalf("Client shutdown failure: %v", err)
		}
		log.Printf("Client shutting down.")
		os.Exit(0)
	}(mshCtx)
}

func initContext(mashupGoodies map[string]interface{}) *sdk.MashupContext {

	var err error
	mashupContext = &sdk.MashupContext{Context: context.Background(), MashupGoodies: mashupGoodies}
	insecure = mashupGoodies["insecure"].(*bool)

	// Initialize local server.
	mashupCertBytes, err = mashupsdk.MashupCert.ReadFile("tls/mashup.crt")
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
	go func() {
		sdk.RegisterMashupServerServer(handshakeServer, &MashupHandshakeServer{})
		handshakeServer.Serve(lis)
	}()

	jsonHandshakeCredentials, err := json.Marshal(handshakeCredentials)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}
	// Setup the onetime use handshake token...
	mashupGoodies["PARAMS"] = append(mashupGoodies["PARAMS"].([]string), "-CREDS="+string(jsonHandshakeCredentials))
	mashupGoodies["PARAMS"] = append(mashupGoodies["PARAMS"].([]string), "-insecure=true")

	// Start mashup..
	err = mashupInit(mashupGoodies)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}

	<-handshakeCompleteChan

	return mashupContext
}

func BootstrapInit(mashupPath string,
	envParams []string,
	params []string,
	insecure *bool) *sdk.MashupContext {

	mashupGoodies := map[string]interface{}{}
	mashupGoodies["MASHUP_PATH"] = mashupPath
	mashupGoodies["ENV"] = envParams
	mashupGoodies["PARAMS"] = params
	mashupGoodies["insecure"] = insecure

	return initContext(mashupGoodies)
}
