package client

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"syscall"

	"github.com/mrjrieke/nute/mashupsdk"
	sdk "github.com/mrjrieke/nute/mashupsdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client bootstrapping is concerned with establishing connection with
// mashup, handshaking, and establishing credential sets.  It
// also sets up signal handling in event of either system
// shutting down.
var handshakeConnectionConfigs *sdk.MashupConnectionConfigs
var clientConnectionConfigs *sdk.MashupConnectionConfigs
var serverConnectionConfigs *sdk.MashupConnectionConfigs

var mashupContext *sdk.MashupContext
var insecure *bool

var handshakeCompleteChan chan bool

var mashupCertBytes []byte

// forkMashup -- starts mashup
func forkMashup(mashupGoodies map[string]interface{}) error {
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
	mashupPath, lookupErr := exec.LookPath(mashupGoodies["MASHUP_PATH"].(string))
	if lookupErr != nil {
		log.Fatalf("Couldn't exec mashup: %v", lookupErr)
	}

	var pid, forkErr = syscall.ForkExec(mashupPath, params, &procAttr)
	if forkErr != nil {
		log.Fatalf("Couldn't exec mashup: %v", forkErr)
	}
	log.Println("Spawned proc", pid)
	mashupGoodies["PID"] = pid

	return forkErr
}

func initContext(mashupApiHandler mashupsdk.MashupApiHandler,
	mashupGoodies map[string]interface{}) *sdk.MashupContext {
	log.Printf("Initializing Mashup.\n")

	handshakeCompleteChan = make(chan bool)
	var err error
	mashupContext = &sdk.MashupContext{Context: context.Background(), MashupGoodies: mashupGoodies}
	insecure = mashupGoodies["insecure"].(*bool)

	// Initialize local server.
	mashupCertBytes, err = sdk.MashupCert.ReadFile("tls/mashup.crt")
	if err != nil {
		log.Fatalf("Couldn't load cert: %v", err)
	}

	mashupKeyBytes, err := sdk.MashupKey.ReadFile("tls/mashup.key")
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
	handshakeConnectionConfigs = &sdk.MashupConnectionConfigs{
		AuthToken: string(hex.EncodeToString([]byte(randomSha256[:]))),
		Port:      int64(lis.Addr().(*net.TCPAddr).Port),
	}

	forkConnectionConfigs := &sdk.MashupConnectionConfigs{
		CallerToken: handshakeConnectionConfigs.AuthToken,
		Port:        handshakeConnectionConfigs.Port,
	}

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	go func() {
		sdk.RegisterMashupServerServer(handshakeServer, &MashupClient{mashupApiHandler: mashupApiHandler})
		handshakeServer.Serve(lis)
	}()

	jsonHandshakeCredentials, err := json.Marshal(forkConnectionConfigs)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}
	// Setup the onetime use handshake token...
	mashupGoodies["PARAMS"] = append(mashupGoodies["PARAMS"].([]string), "-CREDS="+string(jsonHandshakeCredentials))
	mashupGoodies["PARAMS"] = append(mashupGoodies["PARAMS"].([]string), "-insecure=true")

	// Start mashup..
	err = forkMashup(mashupGoodies)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}

	<-handshakeCompleteChan
	log.Printf("Mashup initialized.\n")

	return mashupContext
}

// BootstrapInit - main entry point for bootstrapping the sdk.
// This will fork a mashup, connect with it, and handshake with
// it to establish shared set of credentials to be used in
// future transactions.
func BootstrapInit(mashupPath string,
	mashupApiHandler mashupsdk.MashupApiHandler,
	envParams []string,
	params []string,
	insecure *bool) *sdk.MashupContext {

	mashupGoodies := map[string]interface{}{}
	mashupGoodies["MASHUP_PATH"] = mashupPath
	mashupGoodies["ENV"] = envParams
	mashupGoodies["PARAMS"] = params
	mashupGoodies["insecure"] = insecure

	return initContext(mashupApiHandler, mashupGoodies)
}
