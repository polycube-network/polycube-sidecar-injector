package types

// ServerSettings contains settings about the server
type ServerSettings struct {
	// Port is the port the server will listen to
	Port int
	// CertFile is the server certificate
	CertFile string
	// KeyFile is the private key
	KeyFile string
}
