package utils

import (
	"flag"

	"polycube.network/polycube-sidecar-injector/types"
)

// LoadSettings loads the settings
func LoadSettings() types.ServerSettings {
	servSettings := types.ServerSettings{}

	// Parse command-line flags
	flag.IntVar(&servSettings.Port, "port", 443, "server port.")
	flag.StringVar(&servSettings.CertFile, "certFile", "/etc/mutator/certs/cert.crt", "The TLS certificate")
	flag.StringVar(&servSettings.KeyFile, "keyFile", "/etc/mutator/certs/key.pem", "The TLS private key")
	flag.Parse()

	return servSettings
}
