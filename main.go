package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"polycube.network/polycube-sidecar-injector/mutator"
	"polycube.network/polycube-sidecar-injector/utils"
)

func main() {
	//-----------------------------------------
	// Load settings
	//-----------------------------------------

	settings := utils.LoadSettings()

	//-----------------------------------------
	// Start the server
	//-----------------------------------------

	mutator.StartServer(settings)

	//-----------------------------------------
	// Listen on OS shutdown signal
	//-----------------------------------------

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Infoln("Received shutdown signal. Exiting...")
	mutator.StopServer()
}
