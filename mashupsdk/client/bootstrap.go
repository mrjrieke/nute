package client

import (
	"os/exec"
)

func BootstrapInit(exPath string, envParams []string, params []string) ([]byte, error) {
//	cmd := exec.Command( "C:/Program Files/Java/jdk-13.0.1/bin/javaw.exe", "-cp", ".;lib/calc.jar", "Hello", arg[0])
	// newEnv := append(os.Environ(), "FOO=bar") // Add env
	cmd := exec.Command( exPath, params...)
	if len(envParams) > 0 {
	    cmd.Env = envParams
	}
	return cmd.CombinedOutput()
}