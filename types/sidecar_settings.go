package types

// SidecarSettings contains the settings about the sidecar
type SidecarSettings struct {
	// PolycubeImage is the image that must be applied to the polycube sidecar
	PolycubeImage string `yaml:"polycubeImage"`
}
