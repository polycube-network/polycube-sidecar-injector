package main

import (
	"os"
	"os/signal"
	"syscall"

	"polycube.network/polycube-sidecar-injector/types"

	log "github.com/sirupsen/logrus"
	"polycube.network/polycube-sidecar-injector/mutator"
	"polycube.network/polycube-sidecar-injector/utils"
)

func main() {
	//-----------------------------------------
	// Load settings
	//-----------------------------------------

	servSettings, sidecarSettings := utils.LoadSettings()

	//-----------------------------------------
	// Set polycube
	//-----------------------------------------

	types.SetPolycube(sidecarSettings)
	mutator.SetPolycubePatch()
	mutator.MarshalPatchOperations()

	//-----------------------------------------
	// Start the server
	//-----------------------------------------

	mutator.StartServer(servSettings)

	//-----------------------------------------
	// Listen on OS shutdown signal
	//-----------------------------------------

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Infoln("Received shutdown signal. Exiting...")
	mutator.StopServer()
}
