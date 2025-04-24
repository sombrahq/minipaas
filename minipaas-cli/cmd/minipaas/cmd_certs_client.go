package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type CertsClientArgs struct {
	BaseArgs
	CaDir string `arg:"-c,--ca-dir,required" help:"Directory where CA files (ca.pem, ca-key.pem) are located"`
}

func (args *CertsClientArgs) Run() {
	configFile := filepath.Join(args.Env, "minipaas.yaml")
	cfg, configFile, err := loadConfig(args.Env)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to load configuration file: %s", configFile))

	clientCertDir := filepath.Join(args.Env, cfg.Api.Certs)

	err = os.MkdirAll(clientCertDir, 0755)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to create output directory: %s", clientCertDir))

	caCertFile := filepath.Join(args.CaDir, "ca.pem")
	caKeyFile := filepath.Join(args.CaDir, "ca-key.pem")

	clientKeyFile := filepath.Join(clientCertDir, "key.pem")
	err = runCommand([]string{"openssl", "genrsa", "-out", clientKeyFile, "4096"}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", clientKeyFile))

	clientCSRFile := filepath.Join(clientCertDir, "client.csr")
	err = runCommand([]string{"openssl", "req", "-subj", "/CN=client", "-new", "-key", clientKeyFile, "-out", clientCSRFile}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", clientCSRFile))

	extFile := filepath.Join(clientCertDir, "extfile-client.cnf")
	content := "extendedKeyUsage = clientAuth"
	err = os.WriteFile(extFile, []byte(content), 0644)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", extFile))

	clientCertFile := filepath.Join(clientCertDir, "cert.pem")
	err = runCommand([]string{
		"openssl", "x509", "-req", "-days", "365", "-sha256",
		"-in", clientCSRFile, "-CA", caCertFile, "-CAkey", caKeyFile,
		"-CAcreateserial", "-out", clientCertFile, "-extfile", extFile,
	}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", clientCertFile))

	clientCaFile := filepath.Join(clientCertDir, "ca.pem")
	err = runCommand([]string{"cp", caCertFile, clientCaFile}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to copy file: %s", clientCaFile))

	fmt.Printf("✅ Client for %s certificates: %s\n", cfg.Api.Host, args.CaDir)
}
