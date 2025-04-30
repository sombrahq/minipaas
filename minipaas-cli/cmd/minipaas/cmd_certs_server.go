package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type CertsServerArgs struct {
	Verbose   bool   `arg:"-v,--verbose" help:"Verbose output" default:"false"`
	CN        string `arg:"--cn,required" help:"CN to use for the certificates."`
	OutputDir string `arg:"-o,--output" help:"Directory where certificates will be generated" default:".certs"`
}

func (args *CertsServerArgs) Run() {
	fqdn := args.CN

	err := os.MkdirAll(args.OutputDir, 0755)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to create output directory: %s", args.OutputDir))

	serverCaKeyFile := filepath.Join(args.OutputDir, "ca-key.pem")
	err = runCommand([]string{"openssl", "genrsa", "-out", serverCaKeyFile, "4096"}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", serverCaKeyFile))

	serverCaFile := filepath.Join(args.OutputDir, "ca.pem")
	err = runCommand([]string{
		"openssl", "req", "-new", "-x509", "-days", "4096", "-sha256",
		"-subj", fmt.Sprintf("/CN=%s", fqdn),
		"-key", serverCaKeyFile,
		"-out", serverCaFile,
	}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", serverCaFile))

	serverKeyFile := filepath.Join(args.OutputDir, "server-key.pem")
	err = runCommand([]string{"openssl", "genrsa", "-out", serverKeyFile, "4096"}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", serverKeyFile))

	serverCSRFile := filepath.Join(args.OutputDir, "server.csr")
	err = runCommand([]string{
		"openssl", "req", "-new", "-key", serverKeyFile,
		"-out", serverCSRFile,
		"-subj", fmt.Sprintf("/CN=%s", fqdn),
	}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", serverCSRFile))

	extFile := filepath.Join(args.OutputDir, "extfile.cnf")
	content := fmt.Sprintf(`subjectAltName = DNS:%s
extendedKeyUsage = serverAuth`, fqdn)
	err = os.WriteFile(extFile, []byte(content), 0644)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", extFile))

	serverCertFile := filepath.Join(args.OutputDir, "server-cert.pem")
	err = runCommand([]string{
		"openssl", "x509", "-req", "-days", "365", "-sha256",
		"-in", serverCSRFile,
		"-CA", serverCaFile,
		"-CAkey", serverCaKeyFile,
		"-CAcreateserial",
		"-out", serverCertFile,
		"-extfile", extFile,
	}, args.Verbose)
	checkErrorPanic(err, fmt.Sprintf("❌ Fail to generate file: %s", serverCertFile))
	fmt.Printf("✅ Server %s certificates: %s\n", fqdn, args.OutputDir)
}
