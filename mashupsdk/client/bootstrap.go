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
		Env:   append([]string{"DISPLAY=:0.0"}, mashupGoodies["ENV"].([]string)...),
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
	insecure = mashupGoodies["tls-skip-validation"].(*bool)
	var maxMessageLength int = -1
	if mml, mmlOk := mashupGoodies["maxMessageLength"].(int); mmlOk {
		maxMessageLength = mml
	}

	// Initialize local server.
	serverCert, err := tls.X509KeyPair(mashupCertBytes, sdk.MashupKeyBytes)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	creds := credentials.NewServerTLSFromCert(&serverCert)

	var handshakeServer *grpc.Server
	if maxMessageLength > 0 {
		handshakeServer = grpc.NewServer(grpc.MaxRecvMsgSize(maxMessageLength), grpc.MaxSendMsgSize(maxMessageLength), grpc.Creds(creds))
	} else {
		handshakeServer = grpc.NewServer(grpc.Creds(creds))
	}
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
		if maxMessageLength > 0 {
			InitDialOptions(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageLength), grpc.MaxCallSendMsgSize(maxMessageLength)))
		}
		sdk.RegisterMashupServerServer(handshakeServer, &MashupClient{mashupApiHandler: mashupApiHandler})
		handshakeServer.Serve(lis)
	}()

	jsonHandshakeCredentials, err := json.Marshal(forkConnectionConfigs)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}
	// Setup the onetime use handshake token...
	mashupGoodies["PARAMS"] = append(mashupGoodies["PARAMS"].([]string), "-CREDS="+string(jsonHandshakeCredentials))
	mashupGoodies["PARAMS"] = append(mashupGoodies["PARAMS"].([]string), "-tls-skip-validation=true")

	// Start mashup..
	err = forkMashup(mashupGoodies)
	if err != nil {
		log.Fatalf("Failure to launch: %v", err)
	}

	<-handshakeCompleteChan
	log.Printf("Mashup initialized.\n")

	return mashupContext
}
func BootstrapInit(mashupPath string,
	mashupApiHandler mashupsdk.MashupApiHandler,
	envParams []string,
	params []string,
	insecure *bool) *sdk.MashupContext {
	return BootstrapInitWithMessageExt(mashupPath, mashupApiHandler, envParams, params, insecure, -1)
}

// BootstrapInitWithMessageExt - main entry point for bootstrapping the sdk.
// This will fork a mashup, connect with it, and handshake with
// it to establish shared set of credentials to be used in
// future transactions.
func BootstrapInitWithMessageExt(mashupPath string,
	mashupApiHandler mashupsdk.MashupApiHandler,
	envParams []string,
	params []string,
	insecure *bool, maxMessageLength int) *sdk.MashupContext {

	mashupGoodies := map[string]interface{}{}
	mashupGoodies["MASHUP_PATH"] = mashupPath
	if envParams == nil {
		envParams = []string{}
	}
	mashupGoodies["ENV"] = envParams
	mashupGoodies["PARAMS"] = params
	mashupGoodies["tls-skip-validation"] = insecure
	mashupGoodies["maxMessageLength"] = maxMessageLength

	return initContext(mashupApiHandler, mashupGoodies)
}
