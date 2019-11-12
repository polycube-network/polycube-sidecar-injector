package mutator

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"polycube.network/polycube-sidecar-injector/types"
)

var (
	server *http.Server
)

// StartServer starts the server injecting the sidecar
func StartServer(servSettings types.ServerSettings) {
	log.Infoln("Server Starting...")

	//-----------------------------------------
	// Load x509 key pair
	//-----------------------------------------
	pair, err := tls.LoadX509KeyPair(servSettings.CertFile, servSettings.KeyFile)
	if err != nil {
		log.Fatalf("An error occurred while trying to load key pair: %v", err)
	}

	//-----------------------------------------
	// Set up the server
	//-----------------------------------------
	server = &http.Server{
		Addr: fmt.Sprintf(":%v", servSettings.Port),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				pair,
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", serve)
	server.Handler = mux

	//-----------------------------------------
	// Actually start the server
	//-----------------------------------------
	go func() {
		log.Infof("Started listening on port %d, certificate file in %s, key file in %s", servSettings.Port, servSettings.CertFile, servSettings.KeyFile)

		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("An error occurred while starting the server: %v", err)
		}
	}()
}

// StopServer stops the server
func StopServer() {
	server.Shutdown(context.Background())
	log.Infoln("Good bye!")
}
