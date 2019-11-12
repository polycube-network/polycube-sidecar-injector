package utils

import (
	"flag"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"polycube.network/polycube-sidecar-injector/types"
)

// LoadSettings loads the settings
func LoadSettings() (types.ServerSettings, types.SidecarSettings) {
	servSettings := types.ServerSettings{}
	sidecarSettingsPath := ""

	// Parse command-line flags
	flag.IntVar(&servSettings.Port, "port", 443, "server port.")
	flag.StringVar(&servSettings.CertFile, "certFile", "/etc/mutator/certs/cert.crt", "The TLS certificate")
	flag.StringVar(&servSettings.KeyFile, "keyFile", "/etc/mutator/certs/key.pem", "The TLS private key")
	flag.StringVar(&sidecarSettingsPath, "sideCar", "/etc/mutator/sidecarconfig.yaml", "The sidecar config map")
	flag.Parse()

	// Unmarshall the config yaml file
	sidecarSettings, err := unmarshalSidecarConfig(sidecarSettingsPath)
	if err != nil {
		log.Warnf("An error occurred while parsing the sidecar settings: %s. Going to use the default settings instead.", err)
	}
	return servSettings, sidecarSettings
}

// unmarshalSidecarConfig gets the config from the yaml and returns the struct
// or the default values
func unmarshalSidecarConfig(confFile string) (types.SidecarSettings, error) {
	// Default values
	defaultSidecarSettings := types.SidecarSettings{
		PolycubeImage: types.PolycubeImageLatest,
	}

	// Read the files
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return defaultSidecarSettings, err
	}

	// Actually unmarshal it
	var sidecarSettings types.SidecarSettings
	if err := yaml.Unmarshal(data, &sidecarSettings); err != nil {
		return defaultSidecarSettings, err
	}

	log.Infof("The sidecar settings have been successfully loaded")

	// A little validation here
	if len(sidecarSettings.PolycubeImage) == 0 {
		log.Fatalf("An empty Polycube Image was provided in the sidecar configMap")
	}

	return sidecarSettings, nil
}
